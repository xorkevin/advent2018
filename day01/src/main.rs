use std::error;
use std::fs::File;
use std::io;
use std::io::prelude::*;
use std::str;

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = io::BufReader::new(file);
    let nums = reader
        .lines()
        .filter_map(|i| match i {
            Ok(i) => i.parse::<i32>().ok(),
            Err(_) => None,
        }).collect::<Vec<_>>();
    let sum = nums.iter().sum::<i32>();
    println!("{}", sum);
    Ok(())
}
