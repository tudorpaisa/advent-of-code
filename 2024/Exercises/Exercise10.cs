namespace AdventOfCode.Exercises;

public class Exercise10 : IExercise
{
    public int GetDay() => 10;

    private record Point(int Row, int Col);
    private static readonly List<Point> directions = [
        new(0, 1),
        new(0, -1),
        new(1, 0),
        new(-1, 0),
    ];

    private List<Point> Search(List<List<int>> input, Point point)
    {
        var cell = input[point.Row][point.Col];
        if (cell == 9) return [point];

        input[point.Row][point.Col] = -10;

        List<Point> points = [];

        foreach (var dir in directions)
        {
            var newPoint = new Point(point.Row + dir.Row, point.Col + dir.Col);
            if (newPoint.Row >= input.Count() || newPoint.Col >= input[0].Count() || newPoint.Row < 0 || newPoint.Col < 0 )
            {
                continue;
            }
            var nextCell = input[newPoint.Row][newPoint.Col];
            if (nextCell == -10 || nextCell - cell != 1) continue;
            points.AddRange(Search(input, newPoint));
        }

        input[point.Row][point.Col] = cell;

        return points;
    }
    private List<List<int>> Parse(string[] input)
    {
        List<List<int>> parsedResult = new();
        foreach(var row in input)
        {
            List<int> parsedRow  = new();
            foreach (var c in row)
            {
                int.TryParse(c.ToString(), out var pc);
                parsedRow.Add(pc);
            }
            parsedResult.Add(parsedRow);
        }
        return parsedResult;
    }

    public Result ExecutePart1(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        var parsed = Parse(input);
        var score = 0;

        for (var i = 0; i < parsed.Count(); i++)
        {
            for (var j = 0; j < parsed[0].Count(); j++)
            {
                if (parsed[i][j] == 0)
                {
                    var points = Search(parsed, new(i, j));
                    var strPoints = points.Select(p => $"{p.Row}, {p.Col}").Distinct();
                    score += strPoints.Count();
                }
            }
        }

        return new Result(score, false);
    }

    public Result ExecutePart2(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        var parsed = Parse(input);
        var score = 0;

        for (var i = 0; i < parsed.Count(); i++)
        {
            for (var j = 0; j < parsed[0].Count(); j++)
            {
                if (parsed[i][j] == 0)
                {
                    var points = Search(parsed, new(i, j));
                    score += points.Count();
                }
            }
        }

        return new Result(score, false);
    }
}
