using System.Text.RegularExpressions;
using AdventOfCode.Utilities;

namespace AdventOfCode.Exercises;

public class Exercise3 : IExercise
{
    public int GetDay() => 3;

    private List<string> GetMUlInstructions(string row)
    {
        var results = Regex.Matches(row, @"mul\(\d+,\d+\)");
        return results.Select(m => m.Value).ToList();
    }

    private (int, int) GetDigits(string match)
    {
        var splitString = match.Replace("mul(", "").Replace(")", "").Split(",");
        int.TryParse(splitString.First(), out var first);
        int.TryParse(splitString.Last(), out var second);
        return (first, second);
    }

    private List<string> GetAllInstructions(string row)
    {
        var results = Regex.Matches(row, @"mul\(\d+,\d+\)|do\(\)|don't\(\)");
        return results.Select(m => m.Value).ToList();
    }

    private List<string> FilterInstructions(List<string> instructions)
    {
        List<string> filtered = new();

        var canMul = true;
        foreach (var i in instructions)
        {
            if (i.StartsWith("mul") && canMul) filtered.Add(i);
            else if (i == "do()") canMul = true;
            else if (i == "don't()") canMul = false;
        }

        return filtered;
    }

    public Result ExecutePart1(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);

        var instructions = input.Select(i => GetMUlInstructions(i)).SelectMany(s => s);
        var numberPairs = instructions.Select(i => GetDigits(i));
        var mult = numberPairs.Select(p => p.Item1 * p.Item2);

        return new Result(mult.Sum(), true);
    }


    public Result ExecutePart2(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);

        var instructions = input.Select(i => GetAllInstructions(i)).SelectMany(s => s).ToList();
        var filtered = FilterInstructions(instructions);
        var numberPairs = filtered.Select(i => GetDigits(i));
        var mult = numberPairs.Select(p => p.Item1 * p.Item2);

        return new Result(mult.Sum(), false);
    }

}
