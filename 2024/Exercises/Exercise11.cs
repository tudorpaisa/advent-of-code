namespace AdventOfCode.Exercises;

public class Exercise11 : IExercise
{
    public int GetDay() => 11;

    private List<int> Parse(string input)
    {
        List<int> parsed = new();
        foreach (var raw in input.Split(" "))
        {
            int.TryParse(raw, out var n);
            parsed.Add(n);
        }
        return parsed;
    }

    private void Iterate(List<int> stones)
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

                int.TryParse(leftRaw, out var left);
                int.TryParse(rightRaw, out var right);
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
        var parsed = Parse(input.First());


        return new Result();
    }
}
