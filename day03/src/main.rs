use std::fs;
use std::io;
use std::io::prelude::*;

const PUZZLEINPUT: &str = "input.txt";
const WIDTH: usize = 1000;

fn main() {
    let re = regex::Regex::new(
        r"^#(?P<claimid>\d+) @ (?P<col>\d+),(?P<row>\d+): (?P<width>\d+)x(?P<height>\d+)$",
    )
    .expect("Invalid regex");
    let file = fs::File::open(PUZZLEINPUT).expect("Failed to open file");
    let reader = io::BufReader::new(file);

    let mut claims_count = vec![0; WIDTH * WIDTH];
    let mut claims = vec![0; WIDTH * WIDTH];
    let mut banned_claims = vec![false; 1253];

    for line in reader.lines() {
        let s = line.expect("Failed to read line");
        let caps = re.captures(&s).expect("Invalid input");
        let claimid = caps
            .name("claimid")
            .expect("Invalid input")
            .as_str()
            .parse::<usize>()
            .expect("Failed to parse claimid");
        let col = caps
            .name("col")
            .expect("Invalid input")
            .as_str()
            .parse::<usize>()
            .expect("Failed to parse col");
        let row = caps
            .name("row")
            .expect("Invalid input")
            .as_str()
            .parse::<usize>()
            .expect("Failed to parse row");
        let width = caps
            .name("width")
            .expect("Invalid input")
            .as_str()
            .parse::<usize>()
            .expect("Failed to parse width");
        let height = caps
            .name("height")
            .expect("Invalid input")
            .as_str()
            .parse::<usize>()
            .expect("Failed to parse height");
        for i in row..(row + height) {
            for j in col..(col + width) {
                claims_count[i * WIDTH + j] += 1;
                let mut c2 = &mut claims[i * WIDTH + j];
                if *c2 > 0 {
                    banned_claims[claimid - 1] = true;
                    banned_claims[*c2 - 1] = true;
                }
                *c2 = claimid;
            }
        }
    }

    println!("{}", claims_count.iter().filter(|&i| *i > 1).count());

    for i in banned_claims
        .iter()
        .enumerate()
        .filter(|(_, &i)| !i)
        .map(|(n, _)| n)
    {
        println!("{}", i + 1);
    }
}
