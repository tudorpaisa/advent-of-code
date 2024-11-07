using AdventOfCode.Utilities;

namespace AdventOfCode.Exercises;

public class Exercise2 : IExercise
{
    public int GetDay() => 2;

    public Result ExecutePart1(string inputFile)
    {
        var lines = File.ReadAllLines(inputFile);
        List<int> scores = [];
        foreach (var line in lines)
        {
            var split = line.Split(" ");
            var left = Exercise2Helpers.LeftColumnChoiceMap[split[0]];
            var right = Exercise2Helpers.RightColumnChoiceMap[split[1]];
            var score = ScoreChoices(right, left) + Exercise2Helpers.PointsMap[right];
            scores.Add(score);
        }
        return new Result(scores.Aggregate((a, b) => a + b), true);
    }

    public Result ExecutePart2(string inputFile)
    {
        var lines = File.ReadAllLines(inputFile);
        List<int> scores = [];
        foreach (var line in lines)
        {
            var split = line.Split(" ");
            var left = Exercise2Helpers.LeftColumnChoiceMap[split[0]];
            var right = PickPlayerShape(left, split[1]);
            var score = ScoreChoices(right, left) + Exercise2Helpers.PointsMap[right];
            scores.Add(score);
        }
        return new Result(scores.Aggregate((a, b) => a + b), true);

    }

    private int ScoreChoices(Shape player, Shape opponent)
    {
        // tie
        if (player == opponent) { return 3; }

        // player beats opponent
        if (Exercise2Helpers.WinHierarchy[player] == opponent) { return 6; }
        // opponent beats player
        if (Exercise2Helpers.WinHierarchy[opponent] == player) { return 0; }

        // wtf?!
        throw new ArgumentException($"You're not playing Rock, Paper, Scissors! P={player}, O={opponent}");
    }

    private Shape PickPlayerShape(Shape opponent, string option)
    {
        switch (option.ToUpper())
        {
            case "X":
                // opponent wins
                return Exercise2Helpers.WinHierarchy[opponent];
            case "Z":
                // opponent loses
                return Exercise2Helpers.LoseHierarchy[opponent];
            default:
                // tie
                return opponent;
        }
    }
}

enum Shape
{
    Rock,
    Paper,
    Scissors,
}

static class Exercise2Helpers
{
    public static Dictionary<Shape, int> PointsMap = new()
    {
        {Shape.Rock, 1},
        {Shape.Paper, 2},
        {Shape.Scissors, 3},
    };
    public static Dictionary<string, Shape> LeftColumnChoiceMap = new()
    {
        {"A", Shape.Rock},
        {"B", Shape.Paper},
        {"C", Shape.Scissors},
    };
    public static Dictionary<string, Shape> RightColumnChoiceMap = new()
    {
        {"X", Shape.Rock},
        {"Y", Shape.Paper},
        {"Z", Shape.Scissors},
    };
    public static Dictionary<Shape, Shape> WinHierarchy = new()
    {
        {Shape.Rock, Shape.Scissors},
        {Shape.Paper, Shape.Rock},
        {Shape.Scissors, Shape.Paper},
    };

    public static Dictionary<Shape, Shape> LoseHierarchy = new()
    {
        {Shape.Scissors, Shape.Rock},
        {Shape.Rock, Shape.Paper},
        {Shape.Paper, Shape.Scissors},
    };
}
