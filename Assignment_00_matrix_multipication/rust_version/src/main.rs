use rand::Rng;
use rayon::prelude::*;
const M_SIZE: usize = 2048;

// Mutating every vec in place because copying is expensive

#[derive(Default, Debug)]
pub struct AllArgs {
    a: Vec<u64>,
    b: Vec<u64>,
    c: Vec<u64>,
    d: Vec<u64>,
}

impl AllArgs {
    fn new() -> AllArgs {
        AllArgs {
            a: vec![0; M_SIZE * M_SIZE],
            b: vec![0; M_SIZE * M_SIZE],
            c: vec![0; M_SIZE * M_SIZE],
            d: vec![0; M_SIZE * M_SIZE],
        }
    }
}

fn idx(i: usize, j: usize) -> usize {
    i * M_SIZE + j
}

pub fn fill_rng(u: &mut [u64]) {
    let mut rng = rand::thread_rng();
    for i in 0..M_SIZE {
        for j in 0..M_SIZE {
            u[idx(i, j)] = rng.gen::<u64>() % 1000;
        }
    }
}

pub fn mat_multiply_transposed(u: &[u64], v: &[u64], out: &mut [u64]) {
    for i in 0..M_SIZE {
        for j in 0..M_SIZE {
            for k in 0..M_SIZE {
                out[idx(i, j)] += u[idx(i, k)] * v[idx(j, k)];
            }
        }
    }
}

pub fn transpose(u: &[u64], out: &mut [u64]) {
    for i in 0..M_SIZE {
        for j in 0..M_SIZE {
            out[idx(j, i)] = u[idx(i, j)];
        }
    }
}

pub fn add_matrix(u: &[u64], v: &[u64], out: &mut [u64]) {
    for i in 0..M_SIZE {
        for j in 0..M_SIZE {
            out[idx(i, j)] = u[idx(i, j)] + v[idx(i, j)];
        }
    }
}

fn compute_matrix_combinators_rayon(a: &[u64], b: &[u64]) -> Vec<Vec<u64>> {
    // let mut sum: u64 = 0;
    (0..M_SIZE)
        .into_par_iter()
        .map(|i| {
            let a_row = &a[idx(i, 0)..idx(i, M_SIZE)];
            compute_row_of_sums_rayon(a_row, b)
        })
        .collect()
}

fn compute_row_of_sums_rayon(a_row: &[u64], b: &[u64]) -> Vec<u64> {
    (0..M_SIZE)
        .into_par_iter()
        .map(|j| {
            a_row
                .iter()
                .zip(b[idx(j, 0)..idx(j, M_SIZE)].iter())
                .map(|(x, y)| x * y)
                .sum()
        })
        .collect()
}

pub fn simple_op_transposed(args: &mut AllArgs) {
    add_matrix(&args.c, &args.d, &mut args.a);
    transpose(&args.a, &mut args.c);
    mat_multiply_transposed(&args.b, &args.c, &mut args.a);
}

pub fn simple_op_transposed_vecs(args: &mut AllArgs) -> Vec<Vec<u64>> {
    add_matrix(&args.c, &args.d, &mut args.a);
    transpose(&args.a, &mut args.c);
    // mat_multiply_transposed(&args.b, &args.c, &mut args.a);
    compute_matrix_combinators_rayon(&args.b, &args.c)
}

fn main() {
    let mut args: AllArgs = AllArgs::new();
    fill_rng(&mut args.a);
    fill_rng(&mut args.b);
    fill_rng(&mut args.c);
    fill_rng(&mut args.d);
    // simple_op_transposed(&mut args);
    let _res = simple_op_transposed_vecs(&mut args);
    //println!("{:?}", args.a);
}
