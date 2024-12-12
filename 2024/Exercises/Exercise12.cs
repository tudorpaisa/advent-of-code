namespace AdventOfCode.Exercises;

public class Exercise12 : IExercise
{
    public int GetDay() => 12;

    private List<(int, int)> directions =
    [
        (0,  1),
        (1,  0),
        (-1, 0),
        (0, -1),
    ];
    private class Region
    {
        public char Value { get; set; }
        public List<(int, int)> Tiles { get; set; }
        public long Area { get; set; }
        public long Perimeter { get; set; }
    }

    private void Discover(string[] input, int row, int col, List<(int, int)> coll)
    {
        var currentVar = input[row][col];
        foreach( var d in directions )
        {
            var newRow = row + d.Item1;
            var newCol = col + d.Item2;

            if (newRow >= input.Count() || newCol >= input[0].Count() || newRow < 0 || newCol < 0 || coll.Contains(( newRow, newCol )) || input[newRow][newCol] != currentVar)
            {
                continue;
            }

            coll.Add((newRow, newCol));
            Discover(input, newRow, newCol, coll);
        }
    }

    private int CalculateSides(string[] input, List<(int, int)> coords)
    {
        var count = 0;

        foreach(var coord in coords)
        {
            ( var row, var col ) = coord;
            var val = input[row][col];
            foreach( var d in directions )
            {
                var newRow = row + d.Item1;
                var newCol = col + d.Item2;

                if (newRow >= input.Count() || newCol >= input[0].Count() || newRow < 0 || newCol < 0)
                {
                    count++;
                    continue;
                }

                if (val != input[newRow][newCol]) count++;
            }
        }

        return count;
    }

    private int CalculateCorners(string[] input, List<(int, int)> coords)
    {
        var count = 0;

        Dictionary<string, (int, int)> dirs = new ();
        dirs.Add("top", (-1, 0));
        dirs.Add("bot", (1, 0));
        dirs.Add("left", (0, -1));
        dirs.Add("right", (0, 1));
        dirs.Add("top-left", (-1, -1));
        dirs.Add("top-right", (-1, 1));
        dirs.Add("bot-left", (1, -1));
        dirs.Add("bot-right", (1, 1));

        foreach (var coord in coords)
        {
            ( var row, var col ) = coord;
            var val = input[row][col];

            Dictionary<string, bool> tileSides = dirs.ToDictionary(kvp => kvp.Key, kvp => false);

            foreach(var d in dirs)
            {
                var newRow = row + d.Value.Item1;
                var newCol = col + d.Value.Item2;

                if (newRow >= input.Count() || newCol >= input[0].Count() || newRow < 0 || newCol < 0)
                {
                    tileSides[d.Key] = true;
                    continue;
                }

                if (val != input[newRow][newCol]) tileSides[d.Key] = true;
            }

            if (tileSides["top"] && tileSides["left"])
            {
                count++;
            }
            if (!tileSides["top"] && !tileSides["left"] && tileSides["top-left"])
            {
                count++;
            }
            if (tileSides["top"] && tileSides["right"])
            {
                count++;
            }
            if (!tileSides["top"] && !tileSides["right"] && tileSides["top-right"])
            {
                count++;
            }
            if (tileSides["bot"] && tileSides["left"])
            {
                count++;
            }
            if (!tileSides["bot"] && !tileSides["left"] && tileSides["bot-left"])
            {
                count++;
            }
            if (tileSides["bot"] && tileSides["right"])
            {
                count++;
            }
            if (!tileSides["bot"] && !tileSides["right"] && tileSides["bot-right"])
            {
                count++;
            }
        }


        return count;
    }

    public Result ExecutePart1(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);

        Dictionary<char, List<List<(int, int)>>> regionCollMap = new();

        for (var i = 0; i < input.Count(); i++)
        {
            for (var j = 0; j < input[0].Count(); j++)
            {
                var val = input[i][j];
                if (regionCollMap.ContainsKey(val) && regionCollMap[val].Where(c => c.Contains((i, j))).Count() > 0) continue;

                List<(int, int)> coll = [( i, j )];
                Discover(input, i, j, coll);
                var regionColl = regionCollMap.GetValueOrDefault(val, new());
                regionColl.Add(coll);
                regionCollMap[val] = regionColl;
            }
        }

        List<Region> regions = new();
        foreach(var kvp in regionCollMap)
        {
            foreach( var coll in kvp.Value )
            {
                var region = new Region { Value = input[coll[0].Item1][coll[0].Item2], Tiles = coll, Area = coll.Count() };
                region.Perimeter = CalculateSides(input, coll);
                regions.Add(region);
            }
        }

        return new Result(regions.Select(r => r.Area * r.Perimeter).Sum(), true);
    }

    public Result ExecutePart2(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);

        Dictionary<char, List<List<(int, int)>>> regionCollMap = new();

        for (var i = 0; i < input.Count(); i++)
        {
            for (var j = 0; j < input[0].Count(); j++)
            {
                var val = input[i][j];
                if (regionCollMap.ContainsKey(val) && regionCollMap[val].Where(c => c.Contains((i, j))).Count() > 0) continue;

                List<(int, int)> coll = [( i, j )];
                Discover(input, i, j, coll);
                var regionColl = regionCollMap.GetValueOrDefault(val, new());
                regionColl.Add(coll);
                regionCollMap[val] = regionColl;
            }
        }

        List<Region> regions = new();
        foreach(var kvp in regionCollMap)
        {
            foreach( var coll in kvp.Value )
            {
                var region = new Region { Value = input[coll[0].Item1][coll[0].Item2], Tiles = coll, Area = coll.Count() };
                region.Perimeter = CalculateCorners(input, coll);
                regions.Add(region);
            }
        }

        foreach(var r in regions)
        {
            // Console.WriteLine($"---{r.Value}---");
            // Console.WriteLine($"{r.Tiles.Select(c => string.Join(" | ", string.Join(", ", c) ))}");
            // Console.WriteLine($"{r.Perimeter}, {r.Area}");
        }


        // foreach(var kvp in regionMap)
        // {
        //     Console.WriteLine($"--{kvp.Key}--");
        //     Console.WriteLine($"{string.Join("\n", kvp.Value.Select(c => string.Join(", ", c)))}");
        // }

        return new Result(regions.Select(r => r.Area * r.Perimeter).Sum(), true);
    }
}
