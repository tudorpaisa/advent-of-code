#!/usr/bin/env python3
with open("input1.txt", "r") as f:
    r1, d =  0, dict()
    for card_num, line in enumerate(f.readlines(), 1):
        winning, in_hand = line.split(":")[-1].split("|")
        points = len([i for i in list(map(int, in_hand.split())) if i in list(map(int, winning.split()))])
        r1 += 1 << (points - 1) if points - 1 >= 0 else 0
        d[card_num] = d.get(card_num, 1)
        for i in range(1, points + 1): d[card_num + i] = d.get(card_num + i, 1) + d[card_num]
    print(r1, sum(d.values()))
