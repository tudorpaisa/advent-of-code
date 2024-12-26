namespace AdventOfCode.Exercises;

public class Exercise20 : IExercise
{
    public int GetDay() => 20;

    private List<(int Y, int X)> Directions =
    [
        (0, 1),
        (1, 0),
        (0, -1),
        (-1, 0),
    ];

    private (int Y, int X) FindWhereString(List<List<string>> map, string s)
    {
        for (var i = 0; i < map.Count(); i++)
        {
            for (var j = 0; j < map[0].Count(); j++)
            {
                if (map[i][j] == s) return (i, j);
            }
        }
        return (-1, -1);
    }

    private List<string> DriveableTiles = [".", "S", "E"];

    private List<( int Y, int X )> FindShortestPath(List<List<string>> map, (int Y, int X) start, (int Y, int X) end, List<(int, int)> visited)
    {
        var q = new Queue<((int Y, int X), List<(int, int)>)>();
        q.Enqueue((start, visited));

        List<List<(int Y, int X)>> paths = [];

        while (q.Any())
        {
            var (node, vis) = q.Dequeue();
            vis.Add(node);

            if (node.Y == end.Y && node.X == end.X)
            {
                paths.Add((vis));
            }

            foreach( var dd in Directions )
            {
                var nx = dd.X + node.X;
                var ny = dd.Y + node.Y;

                if (nx < 0 || ny < 0 || nx >= map.First().Count() || ny >= map.Count()) continue;
                if (vis.Contains((ny, nx))) continue;

                // if not drivable at all
                if (!DriveableTiles.Contains(map[ny][nx])) continue;
                // if drivable
                else q.Enqueue(((ny, nx), vis.ToList()));
            }
        }

        return paths.MinBy(p => p.Count()) ?? [];
    }

    private bool InBounds(List<List<string>> map, int y, int x)
    {
        return y < map.Count() && x < map[0].Count() && y >= 0 && x >= 0;
    }

    private List<List<string>> Parse(string[] input)
    {
        return input.Select(row => row.Select(c => c.ToString()).ToList()).ToList();
    }

    private int GetCheats(List<List<string>> map, List<(int Y, int X)> path, int saveAtLeast, int upTo)
    {
        Dictionary<(int Y, int X), int> pathSet = path.Select((p, i) => (p, i)).ToDictionary(p => (p.p.Y, p.p.X), p => p.i);
        var seen = new HashSet<(int, int)>();
        var cheats = 0;
        for (var start = 0; start < path.Count(); start++)
        {
            var (y, x) = path[start];

            for (var end = 2; end <= upTo; end++)
            {
                seen.Clear();
                for (var d = 0; d <= end; d++)
                {
                    var cutCorner = pathSet.GetValueOrDefault((y + d, x + end - d), -1);
                    if (cutCorner != -1 && (cutCorner - start - end) >= saveAtLeast && seen.Add((start, cutCorner)))
                    {
                        cheats++;
                    }

                    cutCorner = pathSet.GetValueOrDefault((y - d, x + end - d), -1);
                    if (cutCorner != -1 && (cutCorner - start - end) >= saveAtLeast && seen.Add((start, cutCorner)))
                    {
                        cheats++;
                    }

                    cutCorner = pathSet.GetValueOrDefault((y + d, x - end + d), -1);
                    if (cutCorner != -1 && (cutCorner - start - end) >= saveAtLeast && seen.Add((start, cutCorner)))
                    {
                        cheats++;
                    }

                    cutCorner = pathSet.GetValueOrDefault((y - d, x - end + d), -1);
                    if (cutCorner != -1 && (cutCorner - start - end) >= saveAtLeast && seen.Add((start, cutCorner)))
                    {
                        cheats++;
                    }
                }
            }
        }
        return cheats;
    }

    public Result ExecutePart1(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        var map = Parse(input);
        var start = FindWhereString(map, "S");
        var end = FindWhereString(map, "E");
        var path = FindShortestPath(map, start, end, new());

        return new Result(GetCheats(map, path, 100, 2), false);
    }

    public Result ExecutePart2(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        var map = Parse(input);
        var start = FindWhereString(map, "S");
        var end = FindWhereString(map, "E");
        var path = FindShortestPath(map, start, end, new());

        return new Result(GetCheats(map, path, 100, 20), false);
    }
}
