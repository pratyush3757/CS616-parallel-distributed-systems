Compilation flags:
```bash
gcc -Wall -Wpedantic -pg -o gemm.out gemm.c
```
Profiler used: `gprof`

### Trying out naive multiplication as a baseline:
Initial Compilation flags:
```bash
gcc -Wall -Wpedantic -pg -o gemm.out gemm.c
```
I started with a naive version which computed the equation as is:  
```
Each sample counts as 0.01 seconds.
  %   cumulative   self              self     total           
 time   seconds   seconds    calls   s/call   s/call  name    
100.18    316.41   316.41        2   158.20   158.20  multiply
  0.01    316.45     0.04        3     0.01     0.01  fill_matrix_random
  0.00    316.46     0.01        1     0.01     0.01  add
  0.00    316.46     0.00        3     0.00     0.00  zero_out_matrix
  0.00    316.46     0.00        2     0.00     0.00  swap_ptrs
  0.00    316.46     0.00        1     0.00   316.42  naive_op
```
This leaves a lot to be optimised, so I started out by simplifying the equation to `A = B(C+D)`.  

### Testing Simplified Equation:
Now we have only one multiplication, giving us a `~2x` speedup:
```
Each sample counts as 0.01 seconds.
  %   cumulative   self              self     total           
 time   seconds   seconds    calls   s/call   s/call  name    
100.14    141.69   141.69        1   141.69   141.69  multiply
  0.04    141.75     0.06        3     0.02     0.02  fill_matrix_random
  0.01    141.76     0.01        1     0.01     0.01  add
  0.01    141.77     0.01                             _init
  0.00    141.77     0.00        2     0.00     0.00  zero_out_matrix
  0.00    141.77     0.00        1     0.00   141.70  simple_eqn_op
  0.00    141.77     0.00        1     0.00     0.00  swap_ptrs
```

### Transposing matrix before multiplication:
As accessing the matrix in a row major order is better for the cache, we transpose the second operand and
change the multiply function accordingly. Our `~7x` speedup is significant here:
```
 Each sample counts as 0.01 seconds.
  %   cumulative   self              self     total           
 time   seconds   seconds    calls   s/call   s/call  name    
 99.57     20.97    20.97        1    20.97    20.97  multiply_transposed
  0.38     21.05     0.08        1     0.08     0.08  transpose
  0.10     21.07     0.02        3     0.01     0.01  fill_matrix_random
  0.10     21.09     0.02        1     0.02     0.02  add
  0.05     21.10     0.01                             _init
  0.00     21.10     0.00        2     0.00     0.00  zero_out_matrix
  0.00     21.10     0.00        1     0.00    21.07  simple_eqn_op_transposed
```

### Changing optimisation flags:
Adding the `Ofast` flag to the compilation options gives a `~3x` speedup:
```
Each sample counts as 0.01 seconds.
  %   cumulative   self              self     total           
 time   seconds   seconds    calls   s/call   s/call  name    
 99.63      7.63     7.63        1     7.63     7.63  simple_eqn_op_transposed
  0.52      7.67     0.04        3     0.01     0.01  fill_matrix_random
```

### Adding Multithreading:
After adding threads and experimenting with the no. of threads, we arrive at a goldilocks zone of 8 threads (`~1.8x` speedup):
```
Each sample counts as 0.01 seconds.
  %   cumulative   self              self     total           
 time   seconds   seconds    calls  ms/call  ms/call  name    
 98.87      3.95     3.95                             multiply_transposed_scoped
  0.75      3.98     0.03        1    30.04    30.04  simple_eqn_op_transposed_multithreaded
  0.50      4.00     0.02        3     6.67     6.67  fill_matrix_random
  0.00      4.00     0.00        1     0.00     0.00  multiply_transposed_multithreaded
```

### Changing arrays to be 1D [2048 * 2048]:
Now, the only dependency free way of optimising the code I could come up was to change the array structure to be 1D.
So `A[2048][2048]` -> `A[2048 * 2048]`. This doesn't give us much, but we are `~10%` faster on average:
```
Each sample counts as 0.01 seconds.
  %   cumulative   self              self     total           
 time   seconds   seconds    calls  ms/call  ms/call  name    
 97.18      3.64     3.64                             multiply_transposed_scoped
  1.60      3.70     0.06        3    20.02    20.02  fill_matrix_random
  1.33      3.75     0.05        1    50.06    50.06  simple_eqn_op_transposed_multithreaded
  0.00      3.75     0.00        1     0.00     0.00  multiply_transposed_multithreaded
```

### Further optimisations:
We can use OpenMP or GSL (Gnu Scientific Library), but they may not statically compile.
Other option can be trying out Rust + Rayon lib, but I have not tested Rust's capabilities here.

### Hardware used for tests:
```
$ lscpu
Architecture:            x86_64
  CPU op-mode(s):        32-bit, 64-bit
  Address sizes:         39 bits physical, 48 bits virtual
  Byte Order:            Little Endian
CPU(s):                  12
  On-line CPU(s) list:   0-11
Vendor ID:               GenuineIntel
  Model name:            Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
    CPU family:          6
    Model:               158
    Thread(s) per core:  2
    Core(s) per socket:  6
    Socket(s):           1
    Stepping:            10
    CPU(s) scaling MHz:  18%
    CPU max MHz:         4500.0000
    CPU min MHz:         800.0000
    BogoMIPS:            5202.65
```