#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <inttypes.h>

uint64_t A[2048][2048], B[2048][2048], C[2048][2048], D[2048][2048];

void fill_matrix_random(uint64_t ARR[2048][2048]) {
    for(size_t i = 0; i < 2048; i++) {
        for(size_t j = 0; j < 2048; j++) {
            ARR[i][j] = rand()%1000;
        }
    }
}
void zero_out_matrix(uint64_t ARR[2048][2048]) {
    memset(ARR, 0, 2048 * 2048 * sizeof(uint64_t));
}

void multiply(uint64_t U[2048][2048], uint64_t V[2048][2048], uint64_t OUT[2048][2048]) {
    zero_out_matrix(OUT);
    for(size_t i = 0; i < 2048; i++) {
        for(size_t j = 0; j < 2048; j++) {
            for(size_t k = 0; k < 2048; k++) {
                OUT[i][j] += U[i][k] * V[k][j];
            }
        }
    }
}

void multiply_transposed(uint64_t U[2048][2048], uint64_t V[2048][2048], uint64_t OUT[2048][2048]) {
    zero_out_matrix(OUT);
    for(size_t i = 0; i < 2048; i++) {
        for(size_t j = 0; j < 2048; j++) {
            for(size_t k = 0; k < 2048; k++) {
                OUT[i][j] += U[i][k] * V[j][k];
            }
        }
    }
}

void add(uint64_t U[2048][2048], uint64_t V[2048][2048], uint64_t OUT[2048][2048]) {
    zero_out_matrix(OUT);
    for(size_t i = 0; i < 2048; i++) {
        for(size_t j = 0; j < 2048; j++) {
            OUT[i][j] = U[i][j] + V[i][j];
        }
    }
}

void transpose(uint64_t U[2048][2048], uint64_t OUT[2048][2048]) {
    for(size_t i = 0; i < 2048; i++) {
        for(size_t j = 0; j < 2048; j++) {
            OUT[j][i] = U[i][j];
        }
    }
}

void swap_ptrs(uint64_t *u, uint64_t *v) {
    uint64_t temp;
    temp = *u;
    *u = *v;
    *v = temp;
}

void naive_op(uint64_t (*a)[2048], uint64_t(*b)[2048], uint64_t(*c)[2048], uint64_t(*d)[2048]) {
    multiply(b, d, a);
    swap_ptrs(&a, &d);
    multiply(b, c, a);
    swap_ptrs(&c, &a);
    add(d, c, a);
}

void simple_eqn_op(uint64_t (*a)[2048], uint64_t(*b)[2048], uint64_t(*c)[2048], uint64_t(*d)[2048]) {
    add(d, c, a);
    swap_ptrs(&a, &d);
    multiply(b, d, a);
}

void naive_op_transposed(uint64_t (*a)[2048], uint64_t(*b)[2048], uint64_t(*c)[2048], uint64_t(*d)[2048]) {
    transpose(d, a);
    multiply_transposed(b, a, d);
    transpose(c, a);
    multiply_transposed(b, a, c);
    add(d, c, a);
}

void simple_eqn_op_transposed(uint64_t (*a)[2048], uint64_t(*b)[2048], uint64_t(*c)[2048], uint64_t(*d)[2048]) {
    add(d, c, a);
    transpose(a, d);
    multiply_transposed(b, d, a);
}

int main(int argc, char *argv[])
{
    srand(time(NULL));
    uint64_t (*a)[2048], (*b)[2048], (*c)[2048], (*d)[2048];
    a = A;
    b = B;
    c = C;
    d = D;
    fill_matrix_random(b);
    fill_matrix_random(c);
    fill_matrix_random(d);
    // printf("%ld", sizeof(a));
    simple_eqn_op_transposed(a, b, c, d);

    // for(size_t i = 0; i < 2048; i++) {
    //     for(size_t j = 0; j < 2048; j++) {
    //         printf("%" PRIu64 ", ", a[i][j]);
    //     }
    //     printf("\n");
    // }
    return 0;
}
