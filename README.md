# konamicode

A dialect of [brainfuck](https://en.wikipedia.org/wiki/Brainfuck) that uses
gamepad button names

| KC     | BF | Meaning |
| -      | -  | -       |
| RIGHT  | >  | Increment data pointer |
| LEFT   | <  | Decrement data pointer |
| UP     | +  | Increment byte at data pointer |
| DOWN   | -  | Decrement byte at data pointer |
| A      | ,  | Accept one byte of input, storing it in the byte at the data pointer |
| B      | .  | Output the byte at the data pointer |
| START  | [  | If the byte at the data pointer is zero, jump to the command after the matching ] |
| SELECT | ]  | If the byte at the data pointer is nonzero, jump to the command after the matching [ |
