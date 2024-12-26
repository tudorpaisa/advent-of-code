namespace AdventOfCode.Exercises;

public class Exercise18 : IExercise
{
    public int GetDay() => 18;

    private int InputSize = 1024;
    private int GridSize = 71;
    private (int Y, int X) Start = (0, 0);
    private (int Y, int X) End => (GridSize - 1, GridSize -1);

    private List<(int Y, int X)> Directions =
    [
        (0, 1),
        (0, -1),
        (1, 0),
        (-1, 0),
    ];

    private List<(int Y, int X)> Parse(string[] input)
    {
        List<(int, int)> res = [];
        foreach(var row in input)
        {
            var p = row.Split(",").Select(int.Parse).ToList();
            res.Add((p[1], p[0]));
        }
        return res;
    }

    private List<List<string>> CreateMap(int size)
    {
        List<List<string>> a = new();
        for (var i = 0;i < size; i++)
        {
            a.Add(Enumerable.Repeat(".", size).ToList());
        }
        return a;
    }

    private void PopulateMap(List<List<string>> map, List<(int Y, int X)> byteLocs, int nBytes)
    {
        for (var i = 0; i < nBytes; i++)
        {
            var (y, x) = byteLocs[i];
            map[y][x] = "#";
        }
    }

    private void PrintMap(List<List<string>> map)
    {
        Console.WriteLine(string.Join("\n", map.Select(r => string.Join("", r))));
    }

    private void PrintMap(List<List<string>> map, int y, int x)
    {
        var nmap = map.ToList().Select(i => i.ToList()).ToList();
        nmap[y][x] = "O";
        Console.WriteLine(string.Join("\n", nmap.Select(r => string.Join("", r))));
    }

    private bool CanBeSolved(List<List<string>> byteMap)
    {
        var pq = new Queue<(int Y, int X, int steps)>();
        pq.Enqueue((Start.Y, Start.X, 0));
        var visited = new HashSet<(int, int)>();

        while (pq.Count > 0)
        {
            var (y, x, steps) = pq.Dequeue();
            if (visited.Contains((y, x))) continue;
            // Console.WriteLine($"{y} {x} {End.Y} {End.X}");

            visited.Add((y, x));

            if (y == End.Y && x == End.X) return true;

            foreach (var (dy, dx) in Directions)
            {
                var ny = y + dy;
                var nx = x + dx;

                if (ny < 0 || nx < 0 || ny >= GridSize || nx >= GridSize) continue;
                if (byteMap[ny][nx] != ".") continue;
                if (visited.Contains((ny, nx))) continue;

                pq.Enqueue((ny, nx, steps + 1));
            }
        }
        return false;
    }

    public Result ExecutePart1(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        var byteMap = CreateMap(GridSize);
        var byteLoc = Parse(input);
        PopulateMap(byteMap, byteLoc, InputSize);

        int nSteps = int.MaxValue;
        var pq = new Queue<(int Y, int X, int steps)>();
        pq.Enqueue((Start.Y, Start.X, 0));
        var visited = new HashSet<(int, int)>();

        while (pq.Count > 0)
        {
            var (y, x, steps) = pq.Dequeue();
            if (visited.Contains((y, x))) continue;

            visited.Add((y, x));

            if (y == End.Y && x == End.X) nSteps = Math.Min(nSteps, steps);

            foreach (var (dy, dx) in Directions)
            {
                var ny = y + dy;
                var nx = x + dx;

                if (ny < 0 || nx < 0 || ny >= GridSize || nx >= GridSize) continue;
                if (byteMap[ny][nx] != ".") continue;
                if (visited.Contains((ny, nx))) continue;

                pq.Enqueue((ny, nx, steps + 1));
            }
        }


        return new Result(nSteps, true);
    }

    public Result ExecutePart2(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        var byteMap = CreateMap(GridSize);
        var byteLoc = Parse(input);
        PopulateMap(byteMap, byteLoc, InputSize);

        var solution = (-1, -1);

        for (var i = InputSize; i < byteLoc.Count; i++)
        {
            var (bY, bX) = byteLoc[i];
            // Console.WriteLine($"{bY} {bX}");
            byteMap[bY][bX] = "#";
            // PrintMap(byteMap);
            if (!CanBeSolved(byteMap))
            {
                solution = (bX, bY);
                break;
            }
        }

        return new Result(solution, true);
    }
}
