namespace AdventOfCode.Exercises;

public class Exercise16 : IExercise
{
    public int GetDay() => 16;

    private enum Direction
    {
        N,
        E,
        S,
        W,
    }

    private Direction GetLeft(Direction d)
    {
        switch (d)
        {
            case Direction.N:
                return Direction.W;
            case Direction.E:
                return Direction.N;
            case Direction.S:
                return Direction.E;
            case Direction.W:
                return Direction.S;
            default:
                break;
        }
        throw new Exception("WTF");
    }

    private Direction GetRight(Direction d)
    {
        switch (d)
        {
            case Direction.N:
                return Direction.E;
            case Direction.E:
                return Direction.S;
            case Direction.S:
                return Direction.W;
            case Direction.W:
                return Direction.N;
            default:
                break;
        }
        throw new Exception("WTF");
    }

    private (int Y, int X) DirectionToCoords(Direction d)
    {
        switch (d)
        {
            case Direction.N:
                return (-1, 0);
            case Direction.E:
                return (0, 1);
            case Direction.S:
                return (1, 0);
            case Direction.W:
                return (0, -1);
            default:
                break;
        }
        throw new Exception("WTF");
    }

    private void PrintMap(string[] map, int y, int x)
    {
        var w = map[0].Count() + 1;
        var strMap = string.Join("\n", map);
        var c = w * y + x;
        strMap = strMap[0..c] + 'X' + strMap[(c+1)..strMap.Count()];
        Console.WriteLine(strMap);
    }

    private (int Y, int X) FindCharacter(string[] input, char c)
    {
        for (var i = 0; i < input.Count(); i++)
        {
            var row = input[i];
            if (row.Contains(c))
            {
                return (i, row.IndexOf(c));
            }
        }
        return (-1, -1);
    }


    private long Run1(string[] map)
    {
        var start = FindCharacter(map, 'S');
        // Console.WriteLine($"{start.Y}, {start.X}");
        var end = FindCharacter(map, 'E');

        map[end.Y] = map[end.Y][0..end.X] + '.' + map[end.Y][(end.X+1)..map[0].Count()];

        var pq = new PriorityQueue<(Direction Direction, long Score, int Y, int X), long>();
        pq.Enqueue((Direction.E, 0, start.Y, start.X), 0);
        var visited = new HashSet<(Direction Direction, int Y, int X)>();

        long minScore = long.MaxValue;

        while (pq.Count > 0)
        {
            var (dir, score, y, x) = pq.Dequeue();
            // Console.WriteLine($"{dir.ToString()} {y} {x}");

            if (y == end.Y && x == end.X)
            {
                minScore = score;
                break;
            }

            if (visited.Contains((dir, y, x))) continue;

            visited.Add((dir, y, x));
            // PrintMap(map, y, x);
            if (map[y][x] == '#') continue;

            var dd = DirectionToCoords(dir);

            var nx = x + dd.X;
            var ny = y + dd.Y;

            if (map[ny][nx] == '.')
            {

                if (!visited.Contains((dir, ny, nx)))
                {
                    pq.Enqueue((dir, score + 1, ny, nx), score + 1);
                }

            }

            Direction leftDir = GetLeft(dir);
            // Console.WriteLine($"{leftDir.ToString()} <- {dir.ToString()}");
            Direction rightDir = GetRight(dir);
            // Console.WriteLine($"{rightDir.ToString()} <- {dir.ToString()}");

            if (!visited.Contains((leftDir, y, x)))
            {
                pq.Enqueue((leftDir, score + 1000, y, x), score + 1000);
            }

            if (!visited.Contains((rightDir, y, x)))
            {
                pq.Enqueue((rightDir, score + 1000, y, x), score + 1000);
            }
        }

        return minScore;
    }

    private class State
    {
        public Direction Direction { get; set; }
        public long Score { get; set; }
        public int Y { get; set; }
        public int X { get; set; }
        public HashSet<(Direction Direction, int Y, int X)> Visited { get; set; }
    }

    private long Run2(string[] map)
    {
        var start = FindCharacter(map, 'S');
        // Console.WriteLine($"{start.Y}, {start.X}");
        var end = FindCharacter(map, 'E');

        // map[start.Y] = map[start.Y][0..start.X] + '.' + map[start.Y][(start.X+1)..map[0].Count()];
        map[end.Y] = map[end.Y][0..end.X] + '.' + map[end.Y][(end.X+1)..map[0].Count()];

        var pq = new PriorityQueue<(Direction Direction, long Score, int Y, int X, HashSet<(int Y, int X)> Path), long>();
        pq.Enqueue((Direction.E, 0, start.Y, start.X, new() {(start.Y, start.X)}), 0);

        var visited = new Dictionary<(Direction Direction, int Y, int X), long>();

        bool CanVisit((Direction, int, int) state, long score)
        {
            if (visited.TryGetValue(state, out var prevScore) && prevScore < score) return false;
            visited[state] = score;
            return true;
        }

        var winningPaths = new List<(int Y, int X)> ();

        long minScore = long.MaxValue;

        while (pq.Count > 0)
        {
            var (dir, score, y, x, path) = pq.Dequeue();
            if (minScore < score) break;
            // Console.WriteLine($"{dir.ToString()} {y} {x} pl={path.Count} wp={winningPaths.Count}");

            if (y == end.Y && x == end.X)
            {
                // Console.WriteLine("HGERE");
                minScore = score;
                winningPaths.AddRange(path);
                continue;
                // break;
            }

            if (!CanVisit((dir, y, x), score)) continue;

            // PrintMap(map, y, x);
            // if (map[y][x] == '#') continue;

            var dd = DirectionToCoords(dir);

            var nx = x + dd.X;
            var ny = y + dd.Y;

            if (map[ny][nx] == '.')
            {
                if (CanVisit((dir, ny, nx), score + 1))
                {
                    var newPath = path.ToHashSet();
                    newPath.Add((ny, nx));
                    pq.Enqueue((dir, score + 1, ny, nx, newPath), score + 1);
                }
            }

            Direction leftDir = GetLeft(dir);
            // Console.WriteLine($"{leftDir.ToString()} <- {dir.ToString()}");
            Direction rightDir = GetRight(dir);
            // Console.WriteLine($"{rightDir.ToString()} <- {dir.ToString()}");

            if (CanVisit((leftDir, y, x), score + 1000))
            {
                pq.Enqueue((leftDir, score + 1000, y, x, path), score + 1000);
            }

            if (CanVisit((rightDir, y, x), score + 1000))
            {
                pq.Enqueue((rightDir, score + 1000, y, x, path), score + 1000);
            }
        }

        return winningPaths.Distinct().Count();
    }

    public Result ExecutePart1(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);

        return new Result(Run1(input), true);
    }

    public Result ExecutePart2(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);

        return new Result(Run2(input), true);
    }
}
