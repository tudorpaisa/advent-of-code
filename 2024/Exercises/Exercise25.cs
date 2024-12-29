namespace AdventOfCode.Exercises;

public class Exercise25 : IExercise
{
    public int GetDay() => 25;

    private (List<List<string>> Locks, List<List<string>> Keys) Parse(string[] input)
    {
        List<List<string>> locks = new();
        List<List<string>> keys = new();

        List<string> acc = new();
        foreach (var row in input)
        {
            if (row == "")
            {
                // if lock
                if (acc.First() == string.Join("", Enumerable.Repeat("#", acc.First().Count())))
                {
                    locks.Add(acc);
                }
                else
                {
                    keys.Add(acc);
                }
                acc = new();
                continue;
            }

            acc.Add(row);
        }
        // if lock
        if (acc.First() == string.Join("", Enumerable.Repeat("#", acc.First().Count())))
        {
            locks.Add(acc);
        }
        else
        {
            keys.Add(acc);
        }

        return (locks, keys);
    }

    private List<int> GetLockHeight(List<string> a)
    {
        List<int> h = Enumerable.Repeat(0, a.First().Count()).ToList();
        foreach (var r in a)
        {
            for (var i = 0; i < r.Count(); i++)
            {
                if (r[i] == '#') h[i]++;
            }
        }
        return h;
    }

    private List<int> GetKeyHeight(List<string> a)
    {
        List<int> h = Enumerable.Repeat(a.Count(), a.First().Count()).ToList();
        foreach (var r in a)
        {
            for (var i = 0; i < r.Count(); i++)
            {
                if (r[i] == '.') h[i]--;
            }
        }
        return h;
    }

    public Result ExecutePart1(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        var (locks, keys) = Parse(input);

        var lockHeights = locks.Select(l => GetLockHeight(l)).ToList();
        var keyHeights = keys.Select(k => GetKeyHeight(k)).ToList();

        var maxHeight = locks.First().Count();

        var count = 0;
        foreach (var lh in lockHeights)
        {
            foreach (var kh in keyHeights)
            {
                var heightSum = lh.Zip(kh, (a, b) => a + b).ToList();
                // Console.WriteLine($"[{string.Join(", ", lh)}] + [{string.Join(", ", kh)}] = [{string.Join(", ", heightSum)}]");
                if (heightSum.All(s => s <= maxHeight)) count++;
            }
        }

        return new(count++, false);
    }

    public Result ExecutePart2(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);

        return new Result();
    }
}
