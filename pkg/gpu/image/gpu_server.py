# ssh stu1047@pilogin.hpc.sjtu.edu.cn
# OVnUMAbopDWK
import time
import paramiko
import os
import requests

def send_back_result(status):
    response_data = {
        'status': status,
        'job_name': job_name,
    }

    try:
        response = requests.post('http://192.168.3.8:6443/gpu_result', json=response_data)
        response.raise_for_status()  # Raise an exception for HTTP errors
    except requests.exceptions.RequestException as e:
        print(f"Error sending data to http://192.168.3.8:6443/gpu_result: {e}")
        response_data['error'] = str(e)


# ================ Configuration ================
job_name = os.getenv('job-name')
partition = os.getenv('partition')
N = os.getenv('N')
ntasks_per_node = os.getenv('ntasks-per-node')
cpus_per_task = os.getenv('cpus-per-task')
gres = os.getenv('gres')

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
#SBATCH --output=result/%j.out
#SBATCH --error=result/%j.err

ulimit -l unlimited
ulimit -s unlimited

module load gcc

./{job_name}
''')

# ================== Connect to the server ==================
hostname = 'pilogin.hpc.sjtu.edu.cn'
port = 22
username = 'stu1047'
password = 'OVnUMAbopDWK'
local_dir = './src'
remote_dir = f'/lustre/home/acct-stu/stu1047/{job_name}'

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

    # upload the job script to the server
    stdin, stdout, stderr = ssh.exec_command(f'sbatch {remote_dir}/{job_name}.slurm')
    output = stdout.read().decode()
    if output == '':
        raise Exception(stderr.read().decode())
    job_id = output.split()[-1]
    print(f"Submitted job {job_id}.")

    # ================== Wait for the result ==================
    while True:
        try:
            sftp.stat(f'/lustre/home/acct-stu/stu1047/result/{job_id}.out')
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
    sftp.get(f'/lustre/home/acct-stu/stu1047/result/{job_id}.out', f'./{job_name}.out')
    print(f"Downloaded result file {job_id}.out.")
    send_back_result('success')

except Exception as e:
    print(f"Exception occurred: {e}")
    send_back_result(f"Exception occurred: {e}")
finally:
    # close the connection
    ssh.close()
