use std::fs::File;
use std::io::prelude::*;

const PUZZLEINPUT: &str = "input.txt";

fn is_reactive(a: &char, b: &char) -> bool {
    a != b && a.to_ascii_uppercase() == b.to_ascii_uppercase()
}

fn reduce(state: &Vec<char>) -> Vec<char> {
    let mut acc = Vec::new();
    for i in state.iter() {
        if {
            match acc.last() {
                Some(x) => is_reactive(i, x),
                None => false,
            }
        } {
            acc.pop();
        } else {
            acc.push(*i);
        }
    }
    acc
}

fn is_not_char(a: &char, b: &char) -> bool {
    a != b && a.to_ascii_uppercase() != b.to_ascii_uppercase()
}

fn main() {
    let file = File::open(PUZZLEINPUT).expect("Failed to open file");
    let init_state = file
        .bytes()
        .map(|x| x.expect("Failed to read byte") as char)
        .filter(char::is_ascii_alphanumeric)
        .collect::<Vec<_>>();

    println!("{}", reduce(&init_state).len());

    let mut max = 99999999;
    for i in ('A' as u8)..=('Z' as u8) {
        let c = i as char;
        let l = reduce(
            &init_state
                .iter()
                .filter(|x| is_not_char(x, &c))
                .map(|&x| x)
                .collect::<Vec<_>>(),
        )
        .len();
        if l < max {
            max = l;
        }
    }
    println!("{}", max);
}
