## 🌌 qrk programming language

Parser implements [**Pratt parsing algorithm**](https://journal.stuffwithstuff.com/2011/03/19/pratt-parsers-expression-parsing-made-easy/)

Virtual machine is [**stack-based**](https://en.wikipedia.org/wiki/Stack_machine)

### 📃 Code snippets

```rs
print("Goodbye universe!");
```

```rs
fn fibanacci(n) {
    if n < 2 {
        return n;
    }
    return fibanacci(n-2) + fibanacci(n-1);
}
fibanacci(42);
```

```rs
fn isAnswerOnTheLifeQuestion(answer) { 
    let expectedAnswer = 42;
    return answer == expectedAnswer;
};
isAnswerOnTheLifeQuestion(69);
```

```rs
let person = {
    "age": 10,
    "live": fn() { print("living..."); },
};
person["live"]();
```

### 🚀 How to run locally

- have **go** installed locally

- install dependecies

```bash
go mod download
```

- launch REPL

```bash
make run
```

- or run code from a file with `.qrk` extension

```bash
make run FILE="example.qrk"
```