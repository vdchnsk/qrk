## 🌌 qrk programming language interpreter

Parser implements [**Pratt parsing algorithm**](https://journal.stuffwithstuff.com/2011/03/19/pratt-parsers-expression-parsing-made-easy/)

### 📃 Code snippets

```rs
print("Goodbye universe!");
```

```rs
let fibanacci = fn(n) {
    if n < 2 {
        return n;
    }
    return fibanacci(n-2) + fibanacci(n-1);
}
fibanacci(42);
```

```rs
let isAnswerOnTheLifeQuestion = fn(answer) { 
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
