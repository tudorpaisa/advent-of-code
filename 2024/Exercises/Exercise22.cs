using AdventOfCode.Utilities;

namespace AdventOfCode.Exercises;

public class Exercise22 : IExercise
{
    public int GetDay() => 22;

    private HashSet<string> ToSet(string s)
    {
        return s.Select(c => c.ToString()).ToHashSet();
    }

    private int GetDigit(long i)
    {
        var s = i.ToString().Select(c => c.ToString()).ToList();
        return int.Parse(s.Last());
    }

    private long NextSecret(long secret)
    {
        secret = secret ^ (secret * 64);
        secret = secret % 16777216;
        secret = secret ^ ((long)Math.Floor((double)secret / 32.0));
        secret = secret % 16777216;
        secret = secret ^ (secret * 2048);
        secret = secret % 16777216;
        return secret;
    }

    private List<long> Parse(string[] input)
    {
        return input.Select(n => long.Parse(n)).ToList();
    }

    public Result ExecutePart1(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        var nums = Parse(input);
        List<long> secrets = [];
        foreach (var n in nums)
        {
            var res = n;
            for (var i = 0; i < 2000; i++)
            {
                res = NextSecret(res);
            }
            secrets.Add(res);
        }

        long s = 0;
        foreach (var n in secrets) s += n;

        return new Result(s, true);
    }

    private class Seller
    {
        public List<long> Secrets { get; set; }
        public List<int> Prices { get; set; }
        public List<int> Diffs { get; set; }
        public Dictionary<(int, int, int, int), int> RangeFirstValues { get; set; }

        public Seller(List<long> secrets, List<int> prices, List<int> diffs)
        {
            Secrets = secrets;
            Prices = prices;
            Diffs = diffs;
            RangeFirstValues = ComputeRangeMaxes(prices, diffs);
        }

        public Dictionary<(int, int, int, int), int> ComputeRangeMaxes(List<int> prices, List<int> diffs)
        {
            Dictionary<(int, int, int, int), int> map = new();

            for (var i = 1; i < diffs.Count() - 4; i++)
            {
                var t = (diffs[i], diffs[i+1], diffs[i+2], diffs[i+3]);
                if (map.ContainsKey(t)) continue;
                map[t] = prices[i+3];
                // Console.WriteLine($"{Secrets[0]} {string.Join(", ", t)} -> {prices[i+3]}");
            }
            return map;
        }
    }

    public Result ExecutePart2(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        var nums = Parse(input);
        List<Seller> sellers = [];
        foreach (var n in nums)
        {
            var res = n;
            List<long> secList = [res];
            List<int> digits = [GetDigit(res)];
            List<int> diffs = [0];
            for (var i = 0; i < 2000; i++)
            {
                res = NextSecret(res);
                var newDigit = GetDigit(res);
                diffs.Add(newDigit - digits.Last());
                digits.Add(newDigit);
            }

            // Console.WriteLine("----");
            // Console.WriteLine(n);
            // diffs.PrintLines();

            sellers.Add(new(secList, digits, diffs));
        }

        Dictionary<(int, int, int, int), int> allRangeSums = new();
        List<(int, int, int, int)> allRanges = sellers.SelectMany(s => s.RangeFirstValues.Keys.ToList()).ToHashSet().ToList();

        foreach (var range in allRanges)
        {
            var rangeSums = 0;
            foreach(var seller in sellers)
            {
                rangeSums += seller.RangeFirstValues.GetValueOrDefault(range, 0);
            }
            allRangeSums[range] = rangeSums;
        }

        foreach(var kvp in allRangeSums)
        {
            // Console.WriteLine($"{string.Join(", ", kvp.Key)} -> {kvp.Value}");
        }

        // foreach(var seller in sellers)
        // {
        //     Console.WriteLine(seller.RangeFirstValues.GetValueOrDefault((-2, 1, -1, 3), -1));
        // }

        return new Result(allRangeSums.Values.Max(), true);
    }
}
