use regex::Regex;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

struct PointList {
    points: Vec<(isize, isize, isize, isize)>,
}

impl PointList {
    fn new() -> Self {
        Self { points: Vec::new() }
    }

    fn add(&mut self, point: (isize, isize, isize, isize)) {
        self.points.push(point);
    }

    fn step(&mut self, t: isize) {
        for &mut (ref mut x, ref mut y, vx, vy) in self.points.iter_mut() {
            *x += t * vx;
            *y += t * vy;
        }
    }

    fn size(&self) -> (usize, usize, isize, isize) {
        let (mut minx, mut miny, _, _) = self.points[0];
        let (mut maxx, mut maxy, _, _) = self.points[0];
        for &(x, y, _, _) in self.points.iter().skip(1) {
            if x < minx {
                minx = x;
            }
            if x > maxx {
                maxx = x;
            }
            if y < miny {
                miny = y;
            }
            if y > maxy {
                maxy = y;
            }
        }
        (
            (maxx - minx + 1) as usize,
            (maxy - miny + 1) as usize,
            minx,
            miny,
        )
    }

    fn print(&self) {
        let (sizex, sizey, minx, miny) = self.size();
        let mut board = vec![vec!['.'; sizex]; sizey];
        for &(x, y, _, _) in self.points.iter() {
            let x = (x - minx) as usize;
            let y = (y - miny) as usize;
            board[y][x] = '#';
        }
        for i in board {
            println!("{}", i.iter().collect::<String>());
        }
        println!("{} {} {} {}", sizex, sizey, minx, miny);
    }
}

fn main() {
    let re = Regex::new(
        r"^.*< ?(?P<posx>-?\d+),  ?(?P<posy>-?\d+)>.*< ?(?P<velx>-?\d+),  ?(?P<vely>-?\d+)>$",
    )
    .expect("Invalid regex");
    let file = File::open(PUZZLEINPUT).expect("Failed to open file");
    let reader = BufReader::new(file);

    let mut points = PointList::new();
    for line in reader.lines() {
        let l = line.expect("Failed to read line");
        let caps = re.captures(&l).expect("Regex cannot capture line");
        let posx = caps
            .name("posx")
            .expect("Posx does not exist")
            .as_str()
            .parse::<isize>()
            .expect("Failed to parse posx");
        let posy = caps
            .name("posy")
            .expect("Posy does not exist")
            .as_str()
            .parse::<isize>()
            .expect("Failed to parse posy");
        let velx = caps
            .name("velx")
            .expect("Velx does not exist")
            .as_str()
            .parse::<isize>()
            .expect("Failed to parse velx");
        let vely = caps
            .name("vely")
            .expect("Vely does not exist")
            .as_str()
            .parse::<isize>()
            .expect("Failed to parse vely");
        points.add((posx, posy, velx, vely));
    }

    let step = 10656;
    points.step(step);
    points.print();
    println!("Step: {}", step);
}
