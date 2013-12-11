# nanomemo

Nanomemo is a command-line implementation of the SuperMemo SM2 algorithm for
Spaced Repetition Study. It presents you with a bunch of flash cards based on a
CSV input, tracks your performance over time for each card, and tries to show
you the right cards at each moment to maximize your effort in forming long-term
memory of the subject matter.

If you've used Anki, it's basically a shitty version of Anki.

## How to queue up facts to learn

You create a CSV file of "<question>,<answer>" lines, then you call nanomemo
like:

```
nanomemo -input=my.csv
```

Additionally, if your questions are URLs or files, you can instruct nanomemo to
open them with `/usr/bin/open` when presenting questions like this:

```
nanomemo -input=my.csv -openqs
```

This is useful for associating names to gravatars :)

## How to use it

Once you've started `nanomemo`, you'll be presented with your question. Try to
think of the answer, then press any key to reveal the answer. `nanomemo` now
expects you to press a key from 0-5, with the following meanings:

* 0: I had no idea.
* 1: Once I saw the answer, I sort of remembered.
* 2: When I saw the answer, I had an "oh yeah..." moment.
* 3: Remembered the answer, but it was difficult.
* 4: Remembered the answer after brief hesitation.
* 5: Immediate recall.

Once you press a key in 0-5, it will move on to the next question. Which number
you push will determine how long it is until you are presented with the question
again.

Running `nanomemo` on the same dataset every day is highly recommended: Old
questions you previously recalled immediately will be presented to you again as
the memories start to fade.
