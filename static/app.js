// --- Game Logic ---

const numRows = 13;
const numCols = 15;
const gameGrid = document.getElementById('game-grid');

// simple lvl : 0 = empty, 1=wall, 2= soft wall
// does not include "safe zone" at start
const level = [];
for (let r = 0; r < numRows; r++) {
    level[r] = [];
    for (let c = 0; c < numCols; c++ ) {
        if ( 
            r === 0 || c === 0 || 
            r === numRows -1 || c === numCols-1 || 
            (r % 2 === 0 && c % 2 === 0)
        ) {
            level[r][c] = 1; // wall
        } else if (Math.random() < 0.3) {
            level[r][c] = 2 // soft wall
        } else {
            level[r][c] = 0; // empty
        }
    }
}

// --- Game State --- 
let player = { row: 1, col: 1 };
let bombs = [];         // {row, col, timer}
let explosions = [];    // {row, col, timer}

// rendering grid
function renderGrid() {
    gameGrid.innerHTML = '';
    for (let r = 0; r < numRows; r++ ) {
        for (let c = 0; c < numCols; c++) {
            const cell = document.createElement('div');
            cell.classList.add('cell');
            if (level[r][c] === 1) cell.classList.add('wall');
            if (level[r][c] === 2) cell.classList.add('softWall');
            if (r === player.row && c === player.col) cell.classList.add('player');
            if (bombs.some( b => b.row === r && b.col === c)) {
                cell.classList.add('bomb');
                cell.textContent = '💣';
            }
            if (explosions.some(e => e.row === r && e.col === c)) {
                cell.classList.add('explosion')
            }
            cell.dataset.row = r;
            cell.dataset.col = c;
            gameGrid.appendChild(cell);
        }
    }
}
// renderGrid();

// event key listender
document.addEventListener('keydown', function(e) {
    e.preventDefault()
    let { row, col } = player;
    // movement
    if (e.key === 'ArrowUp') row--;
    if (e.key === 'ArrowDown') row++;
    if (e.key === 'ArrowLeft') col--;
    if (e.key === 'ArrowRight') col++;
    // move player if possible
    if (
        row >= 0 && row < numRows &&
        col >= 0 && col < numCols &&
        level[row][col] === 0
    ) {
        player.row = row;
        player.col = col;
    }
    // bomb placmeent
    if (e.key === ' ') {
        if (!bombs.some(b => b.row === player.row && b.col === player.col)) {
            const bombPos = { row: player.row, col: player.col, timer: 120 }; // 2s @ 60fps
            bombs.push(bombPos);
        }
    }
});


function getExplosionCells(row, col) {
    const cells = [];
    [[0,0],[0,1],[0,-1],[1,0],[-1,0]].forEach(([dr,dc]) => {
        let r = row + dr, c = col + dc;
        if (r >= 0 && r < numRows && c >= 0 && c < numCols && level[r][c] !== 1) {
            cells.push({ row: r, col: c });
        }
    });
    return cells;
}


function updateGameState() {
    let newExplosions = [];
    // update bombs
    bombs.forEach(bomb => {
        bomb.timer -= 1;
        if (bomb.timer <= 0) {
            const affectedCells = getExplosionCells(bomb.row, bomb.col);
            // remove bomb and create explosion
            affectedCells.forEach(cell =>  {
                newExplosions.push({ row: cell.row, col: cell.col, timer: 25});
                if (level[cell.row][cell.col] === 2) level[cell.row][cell.col] = 0;
            })
        }
    });
    bombs = bombs.filter(bomb => bomb.timer > 0);
    explosions = explosions.concat(newExplosions);

    // update explosions
    explosions.forEach(explosion => {
        explosion.timer -= 1;
        //remove soft wall if needed : logic here
    })
    explosions = explosions.filter(explosion => explosion.timer > 0);
}

function gameLoop() {
    updateGameState(); // update timers, move enemies, handle explosions, etc.
    renderGrid();      // draw everything based on the current state
    requestAnimationFrame(gameLoop);
}

gameLoop();











// --- Unused // old logic ---

// function showExplosion(row,col) {
//     const cells = [];
//     // center and 4 direction ( no wall penetration for now )
//     [ [0,0], [0,1], [0,-1], [1,0], [-1,0] ].forEach(([dr,dc]) => {
//         let r = row + dr, c = col + dc;
//         if (r >= 0 && r < numRows && c >= 0 && c < numCols && level[r][c] !== 1) {
//             cells.push({row:r, col:c});
//             //remove soft wall
//             if (level[r][c] === 2) level[r][c] = 0;
//         }
//     });
//     // show explosion
//     cells.forEach(({row, col}) => {
//         const idx = row * numCols + col;
//         gameGrid.children[idx].classList.add('explosion');
//     });
//     setTimeout( () => {
//         cells.forEach(({row,col}) => {
//             const idx = row * numCols + col;
//             gameGrid.children[idx].classList.remove('explosion');
//         });
//         renderGrid();
//     }, 400);
// }