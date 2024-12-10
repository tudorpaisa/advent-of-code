using System.Text.RegularExpressions;

namespace AdventOfCode.Exercises;

public class Exercise10 : IExercise
{
    public int GetDay() => 10;

    private int ScoreTrailHead(string[] input, int row, int col)
    {
        List<int> scores = [];

        return scores.Sum();
    }

    public Result ExecutePart1(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        var score = 0;
        for (var i = 0; i < input.Count(); i++)
        {

        }

        return new Result(score, false);
    }

    public Result ExecutePart2(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);

        return new Result();
    }
}
