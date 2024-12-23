using AdventOfCode.Utilities;

namespace AdventOfCode.Exercises;

public class Exercise19 : IExercise
{
    private class Node
    {
        public string Value { get; set; }
        public Dictionary<string, Node> Next = new();
    }

    private string END = "*";

    private Dictionary<string, Node> BuildTrie(List<string> words)
    {
        Dictionary<string, Node> dict = new();
        foreach (var w in words)
        {
            var newNode = new Node { Value = w };
            newNode.Next[END] = new();
            dict[w] = newNode;
        }
        return dict;
    }

    private Dictionary<string, Node> BuildTrie(string word, Dictionary<string, Node> dict)
    {
        if (word == "")
        {
            dict[END] = new();
            return dict;
        }

        var s = word[0].ToString();
        Node node = dict.GetValueOrDefault(s, new() { Value = s });
        node.Next = BuildTrie(word[1..], node.Next);
        dict[node.Value] = node;

        return dict;
    }

    private (List<string>, List<string>) Parse(string[] input)
    {
        var dictionary = input.First().Split(", ").ToList();
        var targets = input[2..].ToList();

        return (dictionary, targets);
    }

    private int CountPermutations(string word, Dictionary<string, Node> trie)
    {
        List<int> dp = Enumerable.Repeat(0, word.Count()).ToList();

        for (var i=0; i < word.Count(); i++)
        {
            if (i == 0 || dp[i-1] == 1)
            {
                var t = trie.ToDictionary();
                for (var j = i; j < word.Count(); j++)
                {
                    var s = word[j].ToString();
                    if (!t.ContainsKey(s))
                    {
                        break;
                    }

                    t = t[s].Next;

                    if (t.ContainsKey(END))
                    {
                        dp[j] = 1;
                    }
                }
            }
        }

        Console.WriteLine($"{word} -> {string.Join(" ", dp)}");
        return dp.Sum();
    }

    private bool CanBuild(string word, Dictionary<string, Node> trie, Dictionary<string, Node> t)
    {
        if (word == "") return true;
        List<bool> results = [];
        var s = word[0].ToString();

        if (t.ContainsKey("*"))
        {
            if (CanBuild(word, trie, trie)) return true;
        }

        if (!t.ContainsKey(s))
        {
            return false;
        }
        t = t[s].Next;
        return CanBuild(word[1..], trie, t);
    }

    public int GetDay() => 19;

    public Result ExecutePart1(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        var (dictionary, targets) = Parse(input);
        Dictionary<string, Node> trie = new();
        for (var i = 0; i < dictionary.Count(); i++)
        {
            var w = dictionary[i];
            trie = BuildTrie(w, trie);
        }

        var count = 0;
        for (var i = 0; i < targets.Count(); i++)
        {
            var w = targets[i];
            if (CanBuild(w, trie, trie)) count++;
        }

        return new Result(count, false);
    }

    private Dictionary<string, long> memory = [];
    private List<string> towels = [];

    private long CountArrangements(string word)
    {
        if (memory.TryGetValue(word, out long count)) return count;
        foreach (var towel in towels)
        {
            if (word == "") return 1;
            if (word.StartsWith(towel))
            {
                count += CountArrangements(word[towel.Length..]);
            }
        }
        memory[word] = count;

        return count;
    }

    public Result ExecutePart2(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        var (dictionary, targets) = Parse(input);
        towels = dictionary.ToList();

        long count = 0;
        for (var i = 0; i < targets.Count(); i++)
        {
            var w = targets[i];
            count += CountArrangements(w);
        }

        return new Result(count, false);
    }
}
