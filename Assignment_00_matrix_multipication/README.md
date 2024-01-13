### Compilation flags:
```bash
gcc -Wall -Wpedantic -Ofast -pg -static -o multiply.out multiply.c
```

### NOTE:
The binary attached `multiply.out` has been compiled statically, while all the testes were done using a dynamically linked binary.  
Tested binary compilation flags:
```bash
gcc -Wall -Wpedantic -Ofast -pg -o multiply.out multiply.c
```

### Profiler used:
`gprof`

### Writeup:
The writeup explaining the process is attached as `Report.md` and its pdf version `Report.pdf`.

### Final Times(Static File):
```
Each sample counts as 0.01 seconds.
  %   cumulative   self              self     total           
 time   seconds   seconds    calls  ms/call  ms/call  name    
 90.59      3.56     3.56        7   508.58   508.58  multiply_transposed_scoped
  3.69      3.71     0.15                             random
  1.78      3.78     0.07        3    23.33    23.33  fill_matrix_random
  1.27      3.83     0.05        1    50.00    50.00  simple_eqn_op_transposed_multithreaded
  1.27      3.88     0.05                             random_r
  0.76      3.91     0.03                             rand
  0.51      3.93     0.02                             __memset_avx2_unaligned_erms
  0.13      3.93     0.01                             setstate
  0.00      3.93     0.00        1     0.00   120.00  main
  0.00      3.93     0.00        1     0.00     0.00  multiply_transposed_multithreaded
```