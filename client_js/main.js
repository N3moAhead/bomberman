const { Bomber, PlayerMove } = require('./pkg/bomber/bomber');

/**
 * @typedef {import('./pkg/bomber/message').ClassicStatePayload} ClassicStatePayload
 * @typedef {import('./pkg/bomber/message').PlayerMove} PlayerMove
 */

class Bot {
    /**
     * @param {ClassicStatePayload} state
     * @returns {PlayerMove}
     */
    calcNextMove(state) {
        // Currently a pretty lazy player :(
        return PlayerMove.DO_NOTHING;
    }
}

function main() {
    const newBot = new Bot();
    const b = new Bomber(newBot);
    b.start('ws://localhost:8038/ws');
}

main();
