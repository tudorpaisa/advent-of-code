namespace AdventOfCode.Exercises;

public class Exercise24 : IExercise
{
    public int GetDay() => 24;

    private record Command(string A, string Gate, string B, string Out);

    private (Dictionary<string, int>, List<Command>) Parse(string[] input)
    {
        Dictionary<string, int> wires = [];
        List<Command> commands = [];

        var finishedInit = false;
        foreach (var row in input)
        {
            if (row == "")
            {
                finishedInit = true;
                continue;
            }
            if (finishedInit)
            {
                var s = row.Split(" ");
                // Console.WriteLine(s[4]);
                commands.Add(new(s[0], s[1], s[2], s[4]));
            }
            else
            {
                var s = row.Split(": ");
                wires[s[0]] = int.Parse(s[1]);
            }
        }
        return (wires, commands);
    }

    private Dictionary<string, int> Compute(Dictionary<string, int> wires, List<Command> commands)
    {
        Dictionary<string, int> outWires = wires.ToDictionary();
        List<int> executed = [];
        var zWires = commands.Select(c => c.Out).Where(v => v.StartsWith("z")).ToList();
        var nZWiresCompletes = 0;

        while (nZWiresCompletes != zWires.Count())
        {
            var completedExecutions = 0;

            for (var i = 0; i < commands.Count(); i++)
            {
                var command = commands[i];

                if (executed.Contains(i)) {
                    continue;
                }

                if (!outWires.ContainsKey(command.A) || !outWires.ContainsKey(command.B)) continue;

                var a = outWires[command.A];
                var b = outWires[command.B];
                var c = command.Out;

                if (command.Gate == "AND")
                {
                    if (a == 1 && b == 1) outWires[c] = 1;
                    else outWires[c] = 0;
                }
                else if (command.Gate == "OR")
                {
                    if (a == 1 || b == 1) outWires[c] = 1;
                    else outWires[c] = 0;
                }
                else if (command.Gate == "XOR")
                {
                    if (a != b) outWires[c] = 1;
                    else outWires[c] = 0;
                }
                else
                {
                    throw new Exception("WTF");
                }

                executed.Add(i);
                completedExecutions++;
                if (c.StartsWith("z")) nZWiresCompletes++;
            }

            if (completedExecutions == 0) break;
        }
        return outWires;
    }

    public Result ExecutePart1(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        var (wires, commands) = Parse(input);
        Dictionary<string, int> outWires = Compute(wires, commands);
        var zWireDict = outWires.Where(w => w.Key.StartsWith("z")).ToDictionary();
        var zWireList = zWireDict.Keys.Order().ToList();

        List<int> bits = [];
        foreach (var z in zWireList)
        {
            bits = [zWireDict[z], ..bits];
        }
        long res = 0;
        for (var i = 0; i < bits.Count(); i++)
        {
            var n = bits[i];
            var p = bits.Count() - i - 1;

            res += n * (long)Math.Pow(2, p);
        }

        return new(res, true);
    }

    public Result ExecutePart2(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);

        // this was done from a combination of networkx and excel
        return new Result("bgs,pqc,rjm,swt,wsv,z07,z13,z31", true);
    }
}
