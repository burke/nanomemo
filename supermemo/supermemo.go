/*
package supermemo implements the supermemo SM-2 algorithm, described here, and
documented in more detail at http://www.supermemo.com/english/ol/sm2.htm.

First, there are a few vocabulary items used by this algorithm that we should cover:

EF: Easiness FactMetadataor. Higher means easier. This represents the ease of memorizing a given item.

I(n): Interval, in days, between repetitions of an item. n represents the number of times the item has been seen.

q: The quality of a response, from 0-5, where:

  5 - perfect response;
  4 - correct response after a hesitation;
  3 - correct response recalled with serious difficulty;
  2 - incorrect response where the correct one seemed easy to recall;
  1 - incorrect response, correct one remembered upon seeing answer;
  0 - complete blackout.

With vocabulary out of the way, the algorithm itself is fairly straightforward:

1. Split the knowledge into the smallest possible items. Think flash cards.

2. With each item, associate an initial EF of 2.5.

3. Repeat items using the following intervals, expressed in days:

  I(1)       := 1
  I(2)       := 6
  I(n | n>2) := I(n-1)*EF

4. After each response, assess the quality of the response (q) as described
above in the Vocabulary section.

5. After each response, modify the EF by the formula:

  EF':=EF+(0.1-(5-q)*(0.08+(5-q)*0.02))

6. If the most recent q > 3, reset n; that is, restart repetitions from I(0).

7. After all items are processed, repeat all items where q < 4, until all items
have at least 4.

This algorithm is implemented in fact_metadata.go, on the FactMetadata type.
*/
package supermemo
