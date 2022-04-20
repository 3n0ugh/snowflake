# snowflake

Here is the twitter [blog](https://blog.twitter.com/engineering/en_us/a/2010/announcing-snowflake) about snowflake.

<img width="1126" alt="Screen Shot 2022-04-12 at 10 15 06" src="https://user-images.githubusercontent.com/69458980/162902896-d6982af3-cd83-49de-92b1-0afd080605de.png">

- The first bit is an unused assigned bit.
- The second part consists of a 41-bit timestamp (milliseconds) whose value is
  the offset of the current time relative to a certain time.
- The 5 bits of the third and fourth parts represent the data center and worker node, and the max value is
  <br/> `2^5-1 = 31`.
- The last part consists of 8 bits, which means the length of the serial number generated per
  millisecond per working node, a maximum of `2^8-1 = 4095` IDs can be generated in the same
  millisecond.
- In a distributed environment, a five-bit data center and worker mean that can deploy 31
  data centers. Each data center can deploy up to 31 nodes.
- The binary length of 41 bits is at most `2^41-1 millisecond = 69 years`. So the snowflake
  algorithm can be used for up to 69 years.

## Usage

- Check the [example program](https://github.com/3n0ugh/snowflake/blob/main/example/main.go).
```go
   // Create a node
   n, err := snowflake.NewNode(30, 3)
   if err != nil {
      fmt.Println(err)
   }
    
   // Then, generate a id
   id, err := n.Generate()
   if err != nil {
      fmt.Println(err)
   }

   fmt.Printf("ID: %d\n", id)
   fmt.Printf("String: %s\n", id.String())
   fmt.Printf("Uint64: %d\n", id.UInt64())
   
   fmt.Printf("DecomposeID: %v\n", snowflake.DecomposeID(id))
```

## Test and Benchmarking

- Test:
```shell
  go test -v . 
```
- Benchmark:
```shell
  go test -bench=. -count=10 -benchtime=2s . 
```

## Attention

- If you need to handle IDs in Javascript, use a string instead of uint64. Because
  Javascript's maximum integer value you can safely store 53 bits.
