use std::collections::HashSet;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";
const GRID_START: isize = -1;
const GRID_END: isize = 361;
const DIST_CAP: isize = 10000;

struct Point(isize, isize);

fn distance(a: &Point, b: &Point) -> isize {
    let Point(ax, ay) = a;
    let Point(bx, by) = b;
    (ax - bx).abs() + (ay - by).abs()
}

fn find_closest(p: &Point, points: &Vec<Point>) -> (bool, usize, isize) {
    points.iter().enumerate().skip(1).fold(
        (false, 0, distance(p, &points[0])),
        |(tie, ind, dist), (n, i)| {
            let k = distance(p, i);
            if k < dist {
                (false, n, k)
            } else {
                (k == dist || tie, ind, dist)
            }
        },
    )
}

fn is_edge(p: &Point) -> bool {
    let Point(x, y) = p;
    *x == GRID_START || *y == GRID_START || *x == GRID_END - 1 || *y == GRID_END - 1
}

fn combined_distance(p: &Point, points: &Vec<Point>) -> isize {
    points.iter().fold(0, |acc, i| acc + distance(p, i))
}

fn main() {
    let file = File::open(PUZZLEINPUT).expect("Failed to open file");
    let reader = BufReader::new(file);

    let points = reader
        .lines()
        .map(|x| {
            let k = x
                .expect("Failed to read line")
                .split(", ")
                .map(|x| x.parse::<isize>().expect("Failed to parse num"))
                .collect::<Vec<_>>();
            Point(k[0], k[1])
        })
        .collect::<Vec<_>>();

    let mut in_region = 0;
    let mut edge = HashSet::new();
    let mut counts = vec![0; points.len()];
    for i in GRID_START..GRID_END {
        for j in GRID_START..GRID_END {
            let p = Point(j, i);
            let (tie, ind, _) = find_closest(&p, &points);
            if !tie {
                counts[ind] += 1;
                if is_edge(&p) {
                    edge.insert(ind);
                }
            }
            if combined_distance(&p, &points) < DIST_CAP {
                in_region += 1;
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
            .expect("Cannot find max")
            .1
    );

    println!("{}", in_region);
}
