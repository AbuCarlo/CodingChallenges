# README for Anthony Nassar's reimplementation of the wc Unix command-line utility

The problem is specified here: https://codingchallenges.substack.com/p/coding-challenge-1

## Scope

All functionality described here is implemented: https://man7.org/linux/man-pages/man1/wc.1.html

BSD-style options with two dashes, i.e. `--byte-count`, are not supported, since they're not 
idiomatic to Go.

I have verified that `wc` will not read from `STDIN` if you tell it to read files, to wit:

```bash
[ec2-user@ip-172-31-52-82 ~]$ echo Hello | wc
      1       1       6
[ec2-user@ip-172-31-52-82 ~]$ echo -n Hello | wc
      0       1       5
[ec2-user@ip-172-31-52-82 ~]$ echo Hello | wc -l Art\ of\ War\ -\ English\ -\ UTF-8.txt 
7143 Art of War - English - UTF-8.txt
[ec2-user@ip-172-31-52-82 ~]$ curl --silent https://gutenberg.org/cache/epub/132/pg132.txt > Art\ of\ War\ -\ English\ -\ UTF-8.txt 
[ec2-user@ip-172-31-52-82 ~]$ wc Art\ of\ War\ -\ English\ -\ UTF-8.txt 
  7143  58164 342160 Art of War - English - UTF-8.txt
[ec2-user@ip-172-31-52-82 ~]$ wc -c Art\ of\ War\ -\ English\ -\ UTF-8.txt 
342160 Art of War - English - UTF-8.txt
[ec2-user@ip-172-31-52-82 ~]$ wc -m Art\ of\ War\ -\ English\ -\ UTF-8.txt 
339262 Art of War - English - UTF-8.txt
[ec2-user@ip-172-31-52-82 ~]$ wc -l Art\ of\ War\ -\ English\ -\ UTF-8.txt 
7143 Art of War - English - UTF-8.txt
[ec2-user@ip-172-31-52-82 ~]$ wc -w Art\ of\ War\ -\ English\ -\ UTF-8.txt 
58164 Art of War - English - UTF-8.txt
[ec2-user@ip-172-31-52-82 ~]$ echo -n Hello | wc -c
5
[ec2-user@ip-172-31-52-82 ~]$ echo -n Hello | wc -m
5
[ec2-user@ip-172-31-52-82 ~]$ echo -n Hello | wc -l
0
[ec2-user@ip-172-31-52-82 ~]$ echo -n Hello | wc -w
1
[ec2-user@ip-172-31-52-82 ~]$ echo Hello | wc -c
6
[ec2-user@ip-172-31-52-82 ~]$ echo Hello | wc -m
6
[ec2-user@ip-172-31-52-82 ~]$ echo Hello | wc -w
1
[ec2-user@ip-172-31-52-82 ~]$ echo Hello | wc -l
1
```
 
I downloaded the Chinese version:

```bash
[ec2-user@ip-172-31-52-82 ~]$ curl --silent https://gutenberg.org/cache/epub/23864/pg23864.txt > Art\ of\ War\ -\ Chinese\ -\ UTF-8.txt 
[ec2-user@ip-172-31-52-82 ~]$ echo $LC_CTYPE

[ec2-user@ip-172-31-52-82 ~]$ wc Art\ of\ War\ -\ Chinese\ -\ UTF-8.txt 
  496  3082 42252 Art of War - Chinese - UTF-8.txt
[ec2-user@ip-172-31-52-82 ~]$ LC_CTYPE=en_US.UTF-8 !wc
LC_CTYPE=en_US.UTF-8 wc Art\ of\ War\ -\ Chinese\ -\ UTF-8.txt 
  496  3082 42252 Art of War - Chinese - UTF-8.txt
[ec2-user@ip-172-31-52-82 ~]$ LC_CTYPE=en_US.UTF-8 wc --chars Art\ of\ War\ -\ Chinese\ -\ UTF-8.txt 
27210 Art of War - Chinese - UTF-8.txt
[ec2-user@ip-172-31-52-82 ~]$ LC_CTYPE=en_US.UTF-8 wc --bytes Art\ of\ War\ -\ Chinese\ -\ UTF-8.txt 
42252 Art of War - Chinese - UTF-8.txt
[ec2-user@ip-172-31-52-82 ~]$ 
```

Large books:

```bash
[ec2-user@ip-172-31-52-82 ~]$ curl --silent 'https://gutenberg.org/cache/epub/11894/pg11894.txt' > 'Mahabharata trans. Ganguli.txt'
[ec2-user@ip-172-31-52-82 ~]$ wc Mahabharata\ trans.\ Ganguli.txt 
 14072 154885 913066 Mahabharata trans. Ganguli.txt
 [ec2-user@ip-172-31-52-82 ~]$ curl -s https://gutenberg.org/cache/epub/996/pg996.txt > 'Don Quixote.txt'
[ec2-user@ip-172-31-52-82 ~]$ wc Don\ Quixote.txt 
  43285  430279 2391721 Don Quixote.txt
 ```

 Multiple files:

```bash
[ec2-user@ip-172-31-52-82 ~]$ wc Don\ Quixote.txt Mahabharata\ trans.\ Ganguli.txt Art\ of\ War\ -\ English\ -\ UTF-8.txt 
  43285  430279 2391721 Don Quixote.txt
  14072  154885  913066 Mahabharata trans. Ganguli.txt
   7143   58164  342160 Art of War - English - UTF-8.txt
  64500  643328 3646947 total
```


