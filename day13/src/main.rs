use std::cmp::Ordering;
use std::collections::HashSet;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

enum Turn {
    Left,
    Straight,
    Right,
}

impl Turn {
    fn next(&self) -> Self {
        match self {
            Turn::Left => Turn::Straight,
            Turn::Straight => Turn::Right,
            Turn::Right => Turn::Left,
        }
    }
}

#[derive(Clone, Copy)]
enum Dir {
    Up,
    Right,
    Down,
    Left,
}

impl Dir {
    fn right(&self) -> Self {
        match self {
            Dir::Up => Dir::Right,
            Dir::Right => Dir::Down,
            Dir::Down => Dir::Left,
            Dir::Left => Dir::Up,
        }
    }

    fn left(&self) -> Self {
        match self {
            Dir::Up => Dir::Left,
            Dir::Right => Dir::Up,
            Dir::Down => Dir::Right,
            Dir::Left => Dir::Down,
        }
    }

    fn turn(&self, turn: &Turn) -> Self {
        match turn {
            Turn::Left => self.left(),
            Turn::Straight => *self,
            Turn::Right => self.right(),
        }
    }

    fn cart_char_to_dir(c: char) -> Option<Self> {
        match c {
            '^' => Some(Dir::Up),
            '>' => Some(Dir::Right),
            'v' => Some(Dir::Down),
            '<' => Some(Dir::Left),
            _ => None,
        }
    }
}

struct Cart {
    x: usize,
    y: usize,
    dir: Dir,
    turn: Turn,
}

impl Cart {
    fn new(x: usize, y: usize, dir: Dir) -> Self {
        Self {
            x: x,
            y: y,
            dir: dir,
            turn: Turn::Left,
        }
    }

    fn next_step(&self) -> (usize, usize) {
        match self.dir {
            Dir::Up => (self.x, self.y - 1),
            Dir::Down => (self.x, self.y + 1),
            Dir::Left => (self.x - 1, self.y),
            Dir::Right => (self.x + 1, self.y),
        }
    }

    fn tick(&mut self, board: &Vec<Vec<char>>) {
        match self.dir {
            Dir::Up => self.y -= 1,
            Dir::Down => self.y += 1,
            Dir::Left => self.x -= 1,
            Dir::Right => self.x += 1,
        }
        match board[self.y][self.x] {
            '\\' => {
                self.dir = match self.dir {
                    Dir::Up => Dir::Left,
                    Dir::Down => Dir::Right,
                    Dir::Left => Dir::Up,
                    Dir::Right => Dir::Down,
                }
            }
            '/' => {
                self.dir = match self.dir {
                    Dir::Up => Dir::Right,
                    Dir::Down => Dir::Left,
                    Dir::Left => Dir::Down,
                    Dir::Right => Dir::Up,
                }
            }
            '+' => {
                self.dir = self.dir.turn(&self.turn);
                self.turn = self.turn.next();
            }
            _ => (),
        }
    }
}

impl Eq for Cart {}
impl PartialEq for Cart {
    fn eq(&self, other: &Self) -> bool {
        self.x == other.x && self.y == other.y
    }
}
impl Ord for Cart {
    fn cmp(&self, other: &Self) -> Ordering {
        match self.y.cmp(&other.y) {
            Ordering::Equal => self.x.cmp(&other.x),
            x => x,
        }
    }
}
impl PartialOrd for Cart {
    fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
        Some(self.cmp(other))
    }
}

struct Sim {
    carts: Vec<Cart>,
    board: Vec<Vec<char>>,
}

impl Sim {
    fn new(carts: Vec<Cart>, board: Vec<Vec<char>>) -> Self {
        Self {
            carts: carts,
            board: board,
        }
    }

    fn tick(mut self) -> Self {
        let mut crashed = HashSet::new();
        for n in 0..self.carts.len() {
            if crashed.contains(&n) {
                continue;
            }
            {
                let (x, y) = &self.carts[n].next_step();
                if let Some(k) = self.detect_crash((*x, *y), &crashed) {
                    println!("{} {}", x, y);
                    crashed.insert(n);
                    crashed.insert(k);
                }
            }
            self.carts[n].tick(&self.board);
        }
        let mut next_carts = self
            .carts
            .into_iter()
            .enumerate()
            .filter(|(n, _)| !crashed.contains(n))
            .map(|(_, i)| i)
            .collect::<Vec<_>>();
        next_carts.sort();
        Self {
            carts: next_carts,
            board: self.board,
        }
    }

    fn detect_crash(&self, (x, y): (usize, usize), crashed: &HashSet<usize>) -> Option<usize> {
        for (n, i) in self.carts.iter().enumerate() {
            if crashed.contains(&n) {
                continue;
            }
            if i.x == x && i.y == y {
                return Some(n);
            }
        }
        None
    }
}

fn main() {
    let file = File::open(PUZZLEINPUT).expect("Failed to open file");
    let reader = BufReader::new(file);

    let mut sim = {
        let mut carts = Vec::new();
        let mut board = Vec::new();
        for line in reader.lines() {
            board.push(
                line.expect("Failed to read line")
                    .chars()
                    .collect::<Vec<_>>(),
            );
        }
        for (y, row) in board.iter_mut().enumerate() {
            for (x, i) in row.iter_mut().enumerate() {
                if let Some(k) = Dir::cart_char_to_dir(*i) {
                    match k {
                        Dir::Up => *i = '|',
                        Dir::Down => *i = '|',
                        Dir::Left => *i = '-',
                        Dir::Right => *i = '-',
                    }
                    carts.push(Cart::new(x, y, k));
                }
            }
        }
        Sim::new(carts, board)
    };

    while sim.carts.len() > 1 {
        sim = sim.tick();
    }

    let last_cart = &sim.carts[0];
    println!("{} {}", last_cart.x, last_cart.y);
}
