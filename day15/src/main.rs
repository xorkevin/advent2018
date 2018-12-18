use std::collections::HashMap;
use std::fs::File;
use std::io::prelude::*;
use std::io::BufReader;

const PUZZLEINPUT: &str = "input.txt";

#[derive(Clone, Copy, PartialEq, Eq, Hash)]
struct Pos(usize, usize);

#[derive(Clone)]
struct Entity {
    pos: Pos,
    elf: bool,
    health: isize,
    attack: isize,
}

impl Entity {
    fn new(pos: Pos, elf: bool) -> Self {
        Self {
            pos: pos,
            elf: elf,
            health: 200,
            attack: 3,
        }
    }

    fn dead(&self) -> bool {
        self.health < 1
    }

    fn damage(&self, other: &mut Self) -> bool {
        other.health -= self.attack;
        other.dead()
    }
}

struct Game<'a> {
    elfs: HashMap<Pos, Entity>,
    goblins: HashMap<Pos, Entity>,
    board: &'a Vec<Vec<char>>,
}

impl<'a> Game<'a> {
    fn new(
        elfs: HashMap<Pos, Entity>,
        goblins: HashMap<Pos, Entity>,
        board: &'a Vec<Vec<char>>,
    ) -> Self {
        Self {
            elfs: elfs,
            goblins: goblins,
            board: board,
        }
    }

    fn print(&self) {
        let mut b = self.board.clone();
        for k in self.elfs.keys() {
            b[k.1][k.0] = 'E';
        }
        for k in self.goblins.keys() {
            b[k.1][k.0] = 'G';
        }
        for i in b.iter() {
            println!("{}", i.iter().collect::<String>());
        }
    }
}

fn main() {
    let file = File::open(PUZZLEINPUT).expect("Failed to open file");
    let reader = BufReader::new(file);

    let (elfs, goblins, board) = {
        let mut elfs = HashMap::new();
        let mut goblins = HashMap::new();
        let mut board = Vec::new();
        for (y, line) in reader.lines().enumerate() {
            let mut k = line
                .expect("Failed to read line")
                .chars()
                .collect::<Vec<_>>();
            for (x, i) in k.iter_mut().enumerate() {
                match i {
                    'E' => {
                        elfs.insert(Pos(x, y), Entity::new(Pos(x, y), true));
                        *i = '.';
                    }
                    'G' => {
                        goblins.insert(Pos(x, y), Entity::new(Pos(x, y), false));
                        *i = '.';
                    }
                    _ => (),
                }
            }
            board.push(k);
        }
        (elfs, goblins, board)
    };

    let mut game = Game::new(elfs.clone(), goblins.clone(), &board);
    game.print();
}
