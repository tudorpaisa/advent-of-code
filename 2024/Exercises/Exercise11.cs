namespace AdventOfCode.Exercises;

public class Exercise11 : IExercise
{
    public int GetDay() => 11;

    private List<long> Parse(string input)
    {
        List<long> parsed = new();
        foreach (var raw in input.Split(" "))
        {
            long.TryParse(raw, out var n);
            parsed.Add(n);
        }
        return parsed;
    }

    // Brute-force solution
    private void Iterate(List<long> stones)
    {
        var i = 0;
        var stonesCount = stones.Count();
        while(i < stonesCount)
        {
            if (stones[i] == 0)
            {
                stones[i] = 1;
            }
            else if (stones[i].ToString().Count() % 2 == 0)
            {
                var strNum = stones[i].ToString();
                int mid = strNum.Count() / 2;
                var leftRaw = strNum[0..mid];
                var rightRaw = strNum[mid..strNum.Count()];

                long.TryParse(leftRaw, out var left);
                long.TryParse(rightRaw, out var right);
                stones[i] = right;
                stones.Insert(i, left);

                stonesCount++;
                i++;
            }
            else
            {
                stones[i] = stones[i] * 2024;
            }
            i++;
        }
    }

    // Optimized solution
    private void Iterate(Dictionary<long, long> stoneCount)
    {
        Dictionary<long, long> summation = new();
        var copy = stoneCount.ToDictionary();
        foreach(var kvp in copy)
        {
            var key = kvp.Key;
            var val = kvp.Value;
            if (val == 0) continue;

            if (key == 0)
            {
                var newKey = 1;
                summation[newKey] = summation.GetValueOrDefault(newKey, 0) + val;
                stoneCount[key] = 0;
            }
            else if (key.ToString().Count() % 2 == 0)
            {
                var strNum = key.ToString();
                int mid = strNum.Count() / 2;
                var leftRaw = strNum[0..mid];
                var rightRaw = strNum[mid..strNum.Count()];

                long.TryParse(leftRaw, out var left);
                long.TryParse(rightRaw, out var right);

                summation[left] = summation.GetValueOrDefault(left, 0) + val;
                summation[right] = summation.GetValueOrDefault(right, 0) + val;

                stoneCount[key] = 0;
            }
            else
            {
                var newKey = key * 2024;
                summation[newKey] = summation.GetValueOrDefault(newKey, 0) + val;
                stoneCount[key] = 0;
            }
        }

        foreach (var kvp in summation)
        {
            stoneCount[kvp.Key] = stoneCount.GetValueOrDefault(kvp.Key, 0) + kvp.Value;
        }
    }

    public Result ExecutePart1(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        var stones = Parse(input.First());

        for (var i = 0; i < 25; i++)
        {
            Iterate(stones);
            // Console.WriteLine($"{string.Join(" ", stones)}");
        }

        return new Result(stones.Count(), false);
    }

    public Result ExecutePart2(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        var stones = Parse(input.First());
        Dictionary<long, long> stoneCount = new();

        foreach (var i in stones)
        {
            stoneCount[i] = 1;
        }

        for (var i = 0; i < 75; i++)
        {
            Iterate(stoneCount);
        }

        return new Result(stoneCount.Values.Sum(), false);
    }
}
