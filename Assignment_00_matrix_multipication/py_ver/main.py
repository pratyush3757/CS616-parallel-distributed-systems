import numpy as np

N = 1024
if __name__ == "__main__":
    A = np.random.randn(N,N).astype(np.int64)
    B = np.random.randn(N,N).astype(np.int64)
    C = A @ B
    print(C)

