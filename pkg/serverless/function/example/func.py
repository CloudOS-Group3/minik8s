import numpy as np

# add args
def main(x, y):
    matrix1 = np.array([[1, 2], [3, 4]])
    matrix2 = np.array([[5, 6], [7, 8]])

    # add
    matrix_sum = np.add(matrix1, matrix2)
    print("Add (matrix1 + matrix2):")
    print(matrix_sum)

    # mutiply
    matrix_product = np.dot(matrix1, matrix2)
    print("mutiply (matrix1 * matrix2):")
    print(matrix_product)

    return x + y