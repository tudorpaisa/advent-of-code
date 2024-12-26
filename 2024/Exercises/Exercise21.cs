namespace AdventOfCode.Exercises;

public class Exercise21 : IExercise
{
    public int GetDay() => 21;

    private List<(int Y, int X, char d)> Directions =
    [
        (0, 1, '>'),
        (0, -1, '<'),
        (1, 0, 'v'),
        (-1, 0, '^'),
    ];

    private Dictionary<char, (int Y, int X)> NumberPad = new()
    {
        {'7', (0, 0)}, {'8', (0, 1)}, {'9', (0, 2)},
        {'4', (1, 0)}, {'5', (1, 1)}, {'6', (1, 2)},
        {'1', (2, 0)}, {'2', (2, 1)}, {'3', (2, 2)},
        {'#', (3, 0)}, {'0', (3, 1)}, {'A', (3, 2)},
    };

    private Dictionary<char, (int Y, int X)> DirectionalPad = new()
    {
        {'#', (0, 0)}, {'^', (0, 1)}, {'A', (0, 2)},
        {'<', (1, 0)}, {'v', (1, 1)}, {'>', (1, 2)},
    };

    private string ShortestPath(char key1, char key2, Dictionary<char, (int Y, int X)> pad)
    {
        var (y1, x1) = pad[key1];
        var (y2, x2) = pad[key2];

        var ud = y2 > y1 ? string.Join("", Enumerable.Repeat("v", (y2-y1))) : string.Join("", Enumerable.Repeat("^", (y1-y2)));
        var lr = x2 > x1 ? string.Join("", Enumerable.Repeat(">", (x2-x1))) : string.Join("", Enumerable.Repeat("<", (x1-x2)));

        if (x2 > x1 && (y2, x1) != pad['#']) return $"{ud}{lr}A";
        if ((y1, x2) != pad['#']) return $"{lr}{ud}A";
        return $"{ud}{lr}A";
    }

    private List<string> Sequences(string code, Dictionary<char, (int Y, int X)> pad)
    {
        List<string> keys = [];
        var prevKey = 'A';
        foreach (var c in code)
        {
            keys.Add(ShortestPath(prevKey, c, pad));
            prevKey = c;
        }
        return keys;
    }

    private Dictionary<string, int> SeqCounts(string code)
    {
        var seqs = Sequences(code, DirectionalPad);
        Dictionary<string, int> freqTable = new();
        foreach (var s in seqs)
        {
            freqTable[s] = freqTable.GetValueOrDefault(s, 0) + 1;
        }
        return freqTable;
    }

    private long Complexity(List<string> codes, int nRobots)
    {
        List<Dictionary<string, long>> freqTables = [];
        codes.ForEach(c =>
        {
            freqTables.Add(new() { {string.Join("", Sequences(c, NumberPad)), 1} });
        });

        for (var i = 0; i < nRobots; i++)
        {
            for (var j = 0; j < freqTables.Count(); j++)
            {
                var freqTable = freqTables[j];
                Dictionary<string, long> subTable = new();

                foreach (var kvp in freqTable)
                {
                    foreach(var subCount in SeqCounts(kvp.Key))
                    {
                        subTable[subCount.Key] = subTable.GetValueOrDefault(subCount.Key, 0) + subCount.Value * kvp.Value;
                    }
                }

                freqTables[j] = subTable;
            }
        }

        var complexities = freqTables.Select( t => t.Select(kvp => kvp.Key.Count() * kvp.Value ).Sum() );

        return complexities.Select((c, i) =>
        {
            long.TryParse(codes[i].Replace("A", ""), out var si);
            return c * si;
        }).Sum();
    }

    public Result ExecutePart1(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);

        var comp = 0;

        foreach (var s in input)
        {
            var code2 = Sequences(s, NumberPad);
            var code3 = Sequences(string.Join("", code2), DirectionalPad);
            var code4 = Sequences(string.Join("", code3), DirectionalPad);
            var d = string.Join("", code4.Select(d => string.Join("", d)).ToList());

            int.TryParse(s.Replace("A", ""), out var si);
            comp += si * d.Count();
        }

        return new Result(comp, true);
    }

    public Result ExecutePart2(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);

        return new Result(Complexity(input.ToList(), 25), false);
    }
}
