# ssh stu1047@pilogin.hpc.sjtu.edu.cn
# OVnUMAbopDWK
import time
import paramiko
import os
import requests

result_server_url = 'http://192.168.3.8:6443/gpu_result'


# ================ Configuration ================
job_name = os.getenv('JOB_NAME')
partition = os.getenv('PARTITION')
N = os.getenv('N')
ntasks_per_node = os.getenv('NTASKS_PER_NODE')
cpus_per_task = os.getenv('CPUS_PER_TASK')
gres = os.getenv('GRES')

if not partition:
    partition = 'gpu'
if not N:
    N = 1
if not ntasks_per_node:
    ntasks_per_node = 1
if not cpus_per_task:
    cpus_per_task = 1
if not gres:
    gres = 'gpu:1'

def send_back_result(status):
    response_data = {
        'uuid': job_name,
        'result': status,
        'error': '',
    }

    try:
        response = requests.post(result_server_url, json=response_data)
        response.raise_for_status()  # Raise an exception for HTTP errors
    except requests.exceptions.RequestException as e:
        print(f"Error sending data to {result_server_url}: {e}")
        response_data['error'] = str(e)

local_dir = './src'
remote_dir = f'/lustre/home/acct-stu/stu1047/{job_name}'

# get .cu file name
cu_files = [f for f in os.listdir(local_dir) if f.endswith('.cu')]
if len(cu_files) == 0:
    raise Exception("No .cu file found.")

cuda_file_name = cu_files[0].rstrip('.cu')
# save the job script to a file
job_script_path = f'./{job_name}.slurm'
with open(job_script_path, 'w') as f:
    f.write(f'''#!/bin/bash
#SBATCH --job-name={job_name}
#SBATCH --partition={partition}
#SBATCH -N {N}
#SBATCH --ntasks-per-node={ntasks_per_node}
#SBATCH --cpus-per-task={cpus_per_task}
#SBATCH --gres={gres}
#SBATCH --output=%j.out
#SBATCH --error=%j.err

ulimit -l unlimited
ulimit -s unlimited

module load gcc/11.2.0 cuda/11.8.0

./{cuda_file_name}
''')

# ================== Connect to the server ==================
hostname = 'pilogin.hpc.sjtu.edu.cn'
port = 22
username = 'stu1047'
password = 'OVnUMAbopDWK'

ssh = paramiko.SSHClient()
ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())

try:
    # connect to the server
    ssh.connect(hostname, port, username, password)

    # ================== Submit the job ==================
    sftp = ssh.open_sftp()

    try:
        sftp.mkdir(remote_dir)
        print(f"Created remote directory: {remote_dir}")
    except IOError:
        print(f"Remote directory already exists: {remote_dir}")

    for filename in os.listdir(local_dir):
        local_path = os.path.join(local_dir, filename)
        if os.path.isfile(local_path):
            remote_path = os.path.join(remote_dir, filename)
            sftp.put(local_path, remote_path)

    sftp.put(job_script_path, f'{remote_dir}/{job_name}.slurm')
    sftp.close()

    # compile the cuda file
    # nvcc file.cu -o file -lcublas
    _stdin, _stdout, stderr = ssh.exec_command(f'module load gcc/11.2.0 cuda/11.8.0; nvcc {remote_dir}/{cuda_file_name}.cu -o {remote_dir}/{cuda_file_name} -lcublas')
    output = stderr.read().decode()
    if output != '':
        raise Exception(output)
    print(f"Compiled cuda file {cuda_file_name}.")

    # upload the job script to the server
    stdin, stdout, stderr = ssh.exec_command(f'cd {remote_dir}; sbatch {job_name}.slurm')
    output = stdout.read().decode()
    if output == '':
        raise Exception(stderr.read().decode())
    job_id = output.split()[-1]
    print(f"Submitted job {job_id}.")

    # ================== Wait for the result ==================
    while True:
        try:
            sftp.stat(f'{remote_dir}/{job_id}.out')
            print(f"Result file {job_id}.out is available.")
            break
        except FileNotFoundError:
            print(f"Result file {job_id}.out not found. Waiting...")
            time.sleep(5)
        except Exception as e:
            print(f"Exception occurred: {e}. Reconnecting...")
            ssh.connect(hostname, port, username, password)
            sftp = ssh.open_sftp()

    # download the result file
    sftp.get(f'{remote_dir}/{job_id}.out', f'./{job_name}.out')
    print(f"Downloaded result file {job_id}.out.")
    send_back_result('success')

except Exception as e:
    print(f"Exception occurred: {e}")
    send_back_result(f"Exception occurred: {e}")
finally:
    # close the connection
    ssh.close()
