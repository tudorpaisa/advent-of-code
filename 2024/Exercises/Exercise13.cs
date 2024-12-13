namespace AdventOfCode.Exercises;

public class Exercise13 : IExercise
{
    public int GetDay() => 13;

    private record Point(long X, long Y);
    private record Case(Point A, Point B, Point Target);

    private (int, int) ParseButton(string input, string button)
    {
        var split = input.Replace($"Button {button.ToUpper()}: ", "").Split(", ");
        var xRaw = split.First().Replace("X+", "");
        var yRaw = split.Last().Replace("Y+", "");
        int.TryParse(xRaw, out var x);
        int.TryParse(yRaw, out var y);
        return (x, y);
    }

    private (long, long) ParsePrize(string input, long multiplier)
    {
        var split = input.Replace("Prize: ", "").Split(", ");
        var xRaw = split.First().Replace("X=", "");
        var yRaw = split.Last().Replace("Y=", "");
        long.TryParse(xRaw, out var x);
        long.TryParse(yRaw, out var y);
        return (x + multiplier, y + multiplier);
    }

    private long GetMinSearchBoundary(Point button, Point prize)
    {
        return Math.Min(prize.X / button.X, prize.Y / button.Y);
    }

    private long GetMaxSearchBoundary(Point button, Point prize)
    {
        return Math.Min(prize.X / button.X, prize.Y / button.Y);
    }
    private (long, long) Search(Case c)
    {
        var aMax = GetMaxSearchBoundary(c.A, c.Target);
        var bMax = GetMaxSearchBoundary(c.B, c.Target);

        for (var i = 0; i <= aMax; i++)
        {
            for (var j = 0; j <= bMax; j++)
            {
                var nx = i * c.A.X + j * c.B.X;
                var ny = i * c.A.Y + j * c.B.Y;
                if (nx == c.Target.X && ny == c.Target.Y) return (i, j);
            }
        }

        return (0, 0);
    }

    private (long, long) Cramer(Case c)
    {
        var c1 = c.Target.X;
        var c2 = c.Target.Y;

        var a1 = c.A.X;
        var a2 = c.A.Y;

        var b1 = c.B.X;
        var b2 = c.B.Y;

        var x = ((b2 * c1) - (b1 * c2)) / ((b2 * a1) - (b1 * a2));
        var y = ((a2 * c1) - (a1 * c2)) / ((a2 * b1) - (a1 * b2));

        if (( c1 == (a1 * x + b1 * y) ) && (c2 == (a2 * x + b2 * y))) return (x, y);
        return (0, 0);
    }

    private List<Case> Parse(string[] input, long multiplier)
    {
        List<Case> parsed = new();

        for (var i = 0; i < input.Count(); i += 4)
        {
            var aCoords = ParseButton(input[i], "A");
            var bCoords = ParseButton(input[i+1], "B");
            var prizeCoords = ParsePrize(input[i+2], multiplier);

            var a = new Point(aCoords.Item1, aCoords.Item2);
            var b = new Point(bCoords.Item1, bCoords.Item2);
            var prize = new Point(prizeCoords.Item1, prizeCoords.Item2);
            parsed.Add(new(a, b, prize));
        }
        return parsed;
    }

    public Result ExecutePart1(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        var parsed = Parse(input, 0);

        List<(long, long)> res = new();
        foreach (var i in  parsed)
        {
            res.Add(Cramer(i));
        }

        return new Result(res.Select(i => i.Item1 * 3 + i.Item2 * 1).Sum(), false);
    }

    public Result ExecutePart2(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        var parsed = Parse(input, 10000000000000);

        List<(long, long)> res = new();
        foreach (var i in  parsed)
        {
            res.Add(Cramer(i));
        }

        return new Result(res.Select(i => i.Item1 * 3 + i.Item2 * 1).Sum(), false);
    }
}
