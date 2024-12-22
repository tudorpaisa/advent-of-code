namespace AdventOfCode.Exercises;

public class Exercise17 : IExercise
{
    public int GetDay() => 17;

    private class Setup
    {
        public long A { get; set; }
        public long B { get; set; }
        public long C { get; set; }
        public List<int> Program { get; set; }
        public List<long> Output { get; set; }

        public static Setup Parse(string[] input)
        {
            long.TryParse(input[0].Replace("Register A: ", ""), out var a);
            long.TryParse(input[1].Replace("Register B: ", ""), out var b);
            long.TryParse(input[2].Replace("Register C: ", ""), out var c);

            List<int> program = new();
            foreach (var i in input.Last().Replace("Program: ", "").Split(","))
            {
                int.TryParse(i, out var ip);
                program.Add(ip);
            }
            return new() { A = a, B = b, C = c, Program = program, Output = new() };
        }

        public long InterpretComboOperand(int op)
        {
            switch (op)
            {
                case 0:
                case 1:
                case 2:
                case 3:
                    return op;
                case 4:
                    return A;
                case 5:
                    return B;
                case 6:
                    return C;
                default:
                    break;
            }
            throw new Exception("WTF");
        }

        public int ExecuteInstruction(int ins, int op, int pointer)
        {
            switch (ins)
            {
                // adv
                case 0:
                    A = A / (long)Math.Pow(2, InterpretComboOperand(op));
                    return pointer + 2;
                // bxl
                case 1:
                    B = B ^ op;
                    return pointer + 2;
                // bst
                case 2:
                    B = InterpretComboOperand(op) % 8;
                    return pointer + 2;
                // jnz
                case 3:
                    if (A == 0)
                    {
                        return pointer + 2;
                    }
                    return op;
                // bxc
                case 4:
                    B = B ^ C;
                    return pointer + 2;
                // out
                case 5:
                    Output.Add(InterpretComboOperand(op) % 8);
                    return pointer + 2;
                // bdv
                case 6:
                    B = A / (long)Math.Pow(2, InterpretComboOperand(op));
                    return pointer + 2;
                // cdv
                case 7:
                    C = A / (long)Math.Pow(2, InterpretComboOperand(op));
                    return pointer + 2;
                default:
                    throw new Exception("WWWWWTFFFFFF");
            }
        }

        public void Run()
        {
            int pointer = 0;
            while (pointer < Program.Count)
            {
                var ins = Program[pointer];
                var op = Program[pointer+1];
                pointer = ExecuteInstruction(ins, op, pointer);
            }
        }
    }

    public Result ExecutePart1(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        var setup = Setup.Parse(input);

        setup.Run();

        return new Result(string.Join(",", setup.Output), true);
    }

    public Result ExecutePart2(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        var setup = Setup.Parse(input);

        var oldA = setup.A;
        var oldB = setup.B;
        var oldC = setup.C;

        List<long> validAs = [0];

        var reversed = setup.Program.ToArray().Reverse().ToList();

        for (int i = 0; i < reversed.Count; i++)
        {
            List<long> newAs = [];
            foreach (var a in validAs)
            {
                for (long j = 0; j < 8; j++)
                {
                    setup.A = oldA;
                    setup.B = oldB;
                    setup.C = oldC;
                    setup.Output = new();

                    long newA = a << 3;
                    newA += j;
                    setup.A = newA;
                    setup.Run();
                    if (setup.Output[0] == reversed[i]) newAs.Add(newA);
                }
            }

            validAs = [.. newAs];
        }

        return new Result(validAs.Order().First(), true);
    }
}
