use std::collections::HashSet;
use std::error;
use std::fmt;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";
const GRID_START: isize = -360;
const GRID_END: isize = 720;

#[derive(Debug)]
struct BasicError {
    msg: String,
}

impl BasicError {
    fn new(msg: &str) -> BasicError {
        BasicError {
            msg: msg.to_string(),
        }
    }
}

impl fmt::Display for BasicError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{}", self.msg)
    }
}

impl error::Error for BasicError {
    fn description(&self) -> &str {
        &self.msg
    }

    fn cause(&self) -> Option<&error::Error> {
        None
    }
}

struct Point(isize, isize);

fn distance(a: &Point, b: &Point) -> isize {
    let Point(ax, ay) = a;
    let Point(bx, by) = b;
    (ax - bx).abs() + (ay - by).abs()
}

fn find_closest(p: &Point, points: &Vec<Point>) -> (bool, usize, isize) {
    let mut tie = false;
    let mut ind = 0;
    let mut dist = distance(p, &points[0]);
    for (n, i) in points.iter().enumerate().skip(1) {
        let k = distance(p, i);
        if k < dist {
            tie = false;
            ind = n;
            dist = k;
        } else if k == dist {
            tie = true
        }
    }
    (tie, ind, dist)
}

fn is_edge(p: &Point) -> bool {
    let Point(x, y) = p;
    *x == GRID_START || *y == GRID_START || *x == GRID_END - 1 || *y == GRID_END - 1
}

fn main() -> Result<(), Box<error::Error>> {
    let file = File::open(PUZZLEINPUT)?;
    let reader = BufReader::new(file);

    let points = reader
        .lines()
        .filter_map(|x| x.ok())
        .map(|x| {
            let k = x
                .split(", ")
                .filter_map(|x| x.parse::<isize>().ok())
                .collect::<Vec<_>>();
            Point(k[0], k[1])
        }).collect::<Vec<_>>();

    let mut edge = HashSet::new();
    let mut counts = vec![0; points.len()];
    for i in GRID_START..GRID_END {
        for j in GRID_START..GRID_END {
            let p = Point(j, i);
            let (tie, ind, _) = find_closest(&p, &points);
            if !tie {
                counts[ind] += 1;
            }
            if is_edge(&p) {
                edge.insert(ind);
            }
        }
    }

    println!(
        "{}",
        counts
            .iter()
            .enumerate()
            .filter(|(n, _)| !edge.contains(n))
            .max_by_key(|(_, &x)| x)
            .ok_or(BasicError::new("Cannot find max"))?
            .1
    );

    Ok(())
}