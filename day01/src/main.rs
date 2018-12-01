use std::collections::HashSet;
use std::error::Error;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let mut nums = Vec::new();
    for line in reader.lines() {
        nums.push(line?.parse::<i32>()?);
    }

    let sum = nums.iter().sum::<i32>();
    println!("{}", sum);

    let mut numset = HashSet::new();
    let mut counter = 0;
    numset.insert(counter);
    for i in nums.iter().cycle() {
        counter += i;
        if !numset.insert(counter) {
            break;
        }
    }
    println!("{}", counter);

    Ok(())
}
