#include <pthread.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

#define NTHREADS 8

typedef struct {
    uint64_t (*a);
    uint64_t (*b);
    uint64_t (*c);
    uint64_t (*d);
} matrix_args_t;

typedef struct {
    uint64_t (*U);
    uint64_t (*V);
    uint64_t (*OUT);
    size_t start, rows;
} thread_args_t;

uint64_t A[2048 * 2048], B[2048 * 2048], C[2048 * 2048], D[2048 * 2048];

static inline size_t idx(size_t i, size_t j) {
    return i * 2048 + j;
}

void fill_matrix_random(uint64_t (*ARR)) {
    for(size_t i = 0; i < 2048; i++) {
        for(size_t j = 0; j < 2048; j++) {
            ARR[idx(i,j)] = rand()%1000;
        }
    }
}
void zero_out_matrix(uint64_t (*ARR)) {
    memset(ARR, 0, 2048 * 2048 * sizeof(uint64_t));
}

void multiply(uint64_t (*U), uint64_t (*V), uint64_t (*OUT)) {
    zero_out_matrix(OUT);
    for(size_t i = 0; i < 2048; i++) {
        for(size_t j = 0; j < 2048; j++) {
            for(size_t k = 0; k < 2048; k++) {
                OUT[idx(i,j)] += U[idx(i,k)] * V[idx(k,j)];
            }
        }
    }
}

void multiply_transposed(uint64_t (*U), uint64_t (*V), uint64_t (*OUT)) {
    zero_out_matrix(OUT);
    for(size_t i = 0; i < 2048; i++) {
        for(size_t j = 0; j < 2048; j++) {
            for(size_t k = 0; k < 2048; k++) {
                OUT[idx(i,j)] += U[idx(i,k)] * V[idx(j,k)];
            }
        }
    }
}

void *multiply_scoped(void *x) {
    thread_args_t args;
    args = *(thread_args_t *)x;
    for(size_t i = args.start; i < args.start + args.rows && i < 2048; i++) {
        for(size_t j = 0; j < 2048; j++) {
            for(size_t k = 0; k < 2048; k++) {
                args.OUT[idx(i,j)] += args.U[idx(i,k)] * args.V[idx(k,j)];
            }
        }
    }
    return NULL;
}

void *multiply_transposed_scoped(void *x) {
    thread_args_t args;
    args = *(thread_args_t *)x;
    for(size_t i = args.start; i < args.start + args.rows && i < 2048; i++) {
        for(size_t j = 0; j < 2048; j++) {
            for(size_t k = 0; k < 2048; k++) {
                args.OUT[idx(i,j)] += args.U[idx(i,k)] * args.V[idx(j,k)];
            }
        }
    }
    return NULL;
}

void multiply_multithreaded(uint64_t (*U), uint64_t (*V),
                                   uint64_t (*OUT)) {
    zero_out_matrix(OUT);
    pthread_t threads[NTHREADS];
    thread_args_t args[NTHREADS];
    int status;
    size_t rows = 2048 / NTHREADS;

    for(size_t i = 0; i < NTHREADS; i++) {
        size_t start = rows * i;
        thread_args_t x = {U,V,OUT,start,rows};
        args[i] = x;
        status = pthread_create(&threads[i], NULL,
                                multiply_scoped, (void *) &args[i]);
    }

    for(size_t i = 0; i < NTHREADS; i++) {
        status = pthread_join(threads[i], NULL);
    }
}

void multiply_transposed_multithreaded(uint64_t (*U), uint64_t (*V),
                                   uint64_t (*OUT)) {
    zero_out_matrix(OUT);
    pthread_t threads[NTHREADS];
    thread_args_t args[NTHREADS];
    int status;
    size_t rows = 2048 / NTHREADS;

    for(size_t i = 0; i < NTHREADS; i++) {
        size_t start = rows * i;
        thread_args_t x = {U,V,OUT,start,rows};
        args[i] = x;
        status = pthread_create(&threads[i], NULL,
                                multiply_transposed_scoped, (void *) &args[i]);
    }

    for(size_t i = 0; i < NTHREADS; i++) {
        status = pthread_join(threads[i], NULL);
    }
}

void add(uint64_t (*U), uint64_t (*V), uint64_t (*OUT)) {
    zero_out_matrix(OUT);
    for(size_t i = 0; i < 2048; i++) {
        for(size_t j = 0; j < 2048; j++) {
            OUT[idx(i,j)] = U[idx(i,j)] + V[idx(i,j)];
        }
    }
}

void transpose(uint64_t (*U), uint64_t (*OUT)) {
    for(size_t i = 0; i < 2048; i++) {
        for(size_t j = 0; j < 2048; j++) {
            OUT[idx(j,i)] = U[idx(i,j)];
        }
    }
}

void swap_ptrs(uint64_t (**U), uint64_t (**V)) {
    uint64_t (*temp);
    temp = *U;
    *U = *V;
    *V = temp;
}

void naive_op(matrix_args_t pointers) {
    multiply(pointers.b, pointers.d, pointers.a);
    swap_ptrs(&pointers.a, &pointers.d);
    multiply(pointers.b, pointers.c, pointers.a);
    swap_ptrs(&pointers.c, &pointers.a);
    add(pointers.d, pointers.c, pointers.a);
}

void simple_eqn_op(matrix_args_t pointers) {
    add(pointers.d, pointers.c, pointers.a);
    swap_ptrs(&pointers.a, &pointers.d);
    multiply(pointers.b, pointers.d, pointers.a);
}

void naive_op_transposed(matrix_args_t pointers) {
    transpose(pointers.d, pointers.a);
    multiply_transposed(pointers.b, pointers.a, pointers.d);
    transpose(pointers.c, pointers.a);
    multiply_transposed(pointers.b, pointers.a, pointers.c);
    add(pointers.d, pointers.c, pointers.a);
}

void simple_eqn_op_transposed(matrix_args_t pointers) {
    add(pointers.d, pointers.c, pointers.a);
    transpose(pointers.a, pointers.d);
    multiply_transposed(pointers.b, pointers.d, pointers.a);
}

void naive_op_multithreaded(matrix_args_t pointers) {
    multiply_multithreaded(pointers.b, pointers.d, pointers.a);
    swap_ptrs(&pointers.a, &pointers.d);
    multiply_multithreaded(pointers.b, pointers.c, pointers.a);
    swap_ptrs(&pointers.c, &pointers.a);
    add(pointers.d, pointers.c, pointers.a);
}

void simple_eqn_op_multithreaded(matrix_args_t pointers) {
    add(pointers.d, pointers.c, pointers.a);
    swap_ptrs(&pointers.a, &pointers.d);
    multiply_multithreaded(pointers.b, pointers.d, pointers.a);
}

void naive_op_transposed_multithreaded(matrix_args_t pointers) {
    transpose(pointers.d, pointers.a);
    multiply_transposed_multithreaded(pointers.b, pointers.a, pointers.d);
    transpose(pointers.c, pointers.a);
    multiply_transposed_multithreaded(pointers.b, pointers.a, pointers.c);
    add(pointers.d, pointers.c, pointers.a);
}

void simple_eqn_op_transposed_multithreaded(matrix_args_t pointers) {
    add(pointers.d, pointers.c, pointers.a);
    transpose(pointers.a, pointers.d);
    multiply_transposed_multithreaded(pointers.b, pointers.d, pointers.a);
}

int main(int argc, char *argv[])
{
    srand(time(NULL));
    matrix_args_t pointers;
    pointers.a = A;
    pointers.b = B;
    pointers.c = C;
    pointers.d = D;
    fill_matrix_random(pointers.b);
    fill_matrix_random(pointers.c);
    fill_matrix_random(pointers.d);
    // printf("%ld", sizeof(a));
    simple_eqn_op_transposed_multithreaded(pointers);
    // zero_out_matrix(pointers.a);
    // swap_ptrs(&pointers.a, &pointers.d);
    // for(size_t i = 0; i < 2048; i++) {
    //     for(size_t j = 0; j < 2048; j++) {
    //         printf("%" PRIu64 ", ", pointers.a[i][j]);
    //     }
    //     printf("\n");
    // }
    return 0;
}
