const PUZZLE_INPUT: isize = 9810;
const BOARD_SIZE: usize = 300;

fn power_level(x: isize, y: isize) -> isize {
    let id = x + 10;
    (id * y + PUZZLE_INPUT) * id / 100 % 10 - 5
}

fn main() {
    let board = {
        let mut board = vec![vec![0; BOARD_SIZE]; BOARD_SIZE];
        for (y, row) in board.iter_mut().enumerate() {
            for (x, i) in row.iter_mut().enumerate() {
                *i = power_level((x + 1) as isize, (y + 1) as isize)
            }
        }
        board
    };
    let partial = {
        let mut partial = vec![vec![0; BOARD_SIZE + 1]; BOARD_SIZE + 1];
        for y in 1..partial.len() {
            for x in 1..partial.len() {
                partial[y][x] = board[y - 1][x - 1] + partial[y][x - 1] + partial[y - 1][x]
                    - partial[y - 1][x - 1];
            }
        }
        partial
    };
    {
        let mut mp = 0;
        let mut mx = 0;
        let mut my = 0;
        for y in 0..partial.len() - 3 {
            for x in 0..partial.len() - 3 {
                let power =
                    partial[y + 3][x + 3] - partial[y][x + 3] - partial[y + 3][x] + partial[y][x];
                if power > mp {
                    mp = power;
                    mx = x + 1;
                    my = y + 1;
                }
            }
        }
        println!("{} {} {}", mp, mx, my);
    }
    {
        let mut mp = 0;
        let mut mx = 0;
        let mut my = 0;
        let mut ms = 0;
        for i in 1..partial.len() {
            for y in 0..partial.len() - i {
                for x in 0..partial.len() - i {
                    let power = partial[y + i][x + i] - partial[y][x + i] - partial[y + i][x]
                        + partial[y][x];
                    if power > mp {
                        mp = power;
                        mx = x + 1;
                        my = y + 1;
                        ms = i;
                    }
                }
            }
        }
        println!("{} {} {} {}", mp, mx, my, ms);
    }
}
