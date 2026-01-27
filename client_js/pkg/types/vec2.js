class Vec2 {
    constructor(x, y) {
        this.x = x;
        this.y = y;
    }

    toString() {
        return `Vec2{X: ${this.x}, Y: ${this.y}}`;
    }

    add(other) {
        return new Vec2(this.x + other.x, this.y + other.y);
    }

    sub(other) {
        return new Vec2(this.x - other.x, this.y - other.y);
    }

    mul(scalar) {
        return new Vec2(this.x * scalar, this.y * scalar);
    }

    lengthSq() {
        return this.x * this.x + this.y * this.y;
    }

    dot(other) {
        return this.x * other.x + this.y * other.y;
    }
}

module.exports = { Vec2 };
