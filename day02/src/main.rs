use std::collections::HashMap;
use std::error::Error;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

fn main() -> Result<(), Box<Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let lines = {
        let mut lines = Vec::new();
        for line in reader.lines() {
            lines.push(line?);
        }
        lines
    };

    let (twice, thrice) = {
        let mut twice = 0;
        let mut thrice = 0;
        for line in lines.iter() {
            let (is_twice, is_thrice) = is_twice_or_thrice(line);
            if is_twice {
                twice += 1;
            }
            if is_thrice {
                thrice += 1;
            }
        }
        (twice, thrice)
    };

    println!("{}", twice * thrice);

    for line1 in lines.iter() {
        for line2 in lines.iter() {
            if is_single_diff(line1, line2) {
                println!("{}", remove_single_diff(line1, line2));
                return Ok(());
            }
        }
    }

    Ok(())
}

fn is_twice_or_thrice(line: &str) -> (bool, bool) {
    let seen = {
        let mut seen = HashMap::new();
        for c in line.chars() {
            *seen.entry(c).or_insert(0) += 1;
        }
        seen
    };

    let mut twice = false;
    let mut thrice = false;
    for (_, &v) in seen.iter() {
        twice = twice || v == 2;
        thrice = thrice || v == 3;
    }

    (twice, thrice)
}

fn is_single_diff(line1: &str, line2: &str) -> bool {
    if line1.len() != line2.len() {
        return false;
    }

    let mut diff = false;
    for (c1, c2) in line1.chars().zip(line2.chars()) {
        if c1 != c2 {
            if !diff {
                diff = true;
            } else {
                return false;
            }
        }
    }
    diff
}

fn remove_single_diff(line1: &str, line2: &str) -> String {
    for (i, (c1, c2)) in line1.chars().zip(line2.chars()).enumerate() {
        if c1 != c2 {
            let mut s = String::from(line1);
            s.remove(i);
            return s;
        }
    }
    String::new()
}
