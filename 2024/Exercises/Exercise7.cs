namespace AdventOfCode.Exercises;

public class Exercise7 : IExercise
{
    public int GetDay() => 7;

    public Result ExecutePart1(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        List<long> evaluatable = new();
        foreach (var row in input)
        {
            (var lhs, var rhs) = ParseRow(row);
            var eval = CanEvaluate(lhs, rhs, "+", new string[] { "+", "*" }, 0);
            if (eval) evaluatable.Add(lhs);
        }

        return new Result(evaluatable.Sum(), true);
    }

    private (long, List<long>) ParseRow(string row)
    {
        var split1 = row.Split(": ");
        long.TryParse(split1.First(), out var lhs);
        List<long> rhs = new();
        foreach (var rawN in split1.Last().Split(" "))
        {
            int.TryParse(rawN, out var n);
            rhs.Add(n);
        }

        return (lhs, rhs);

    }

    private bool CanEvaluate(long target, List<long> numbers, string sign, string[] signs, long accumulator)
    {
        if (numbers.Count() == 0) return accumulator == target;

        var newNumbers = numbers.ToList();
        var nextNumber = newNumbers.First();
        newNumbers.RemoveAt(0);

        var newAccumulator = accumulator;

        switch (sign)
        {
            case "+":
                newAccumulator += nextNumber;
                break;
            case "*":
                newAccumulator *= nextNumber;
                break;
            case "||":
                var accStr = newAccumulator.ToString();
                var nexNumStr = nextNumber.ToString();
                var newAccStr = accStr + nexNumStr;
                long.TryParse(newAccStr, out newAccumulator);
                break;
            default:
                break;
        }

        foreach (var s in signs)
        {
            if (CanEvaluate(target, newNumbers, s, signs,  newAccumulator)) return true;
        }
        return false;
    }

    public Result ExecutePart2(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        List<long> evaluatable = new();
        foreach (var row in input)
        {
            (var lhs, var rhs) = ParseRow(row);
            var eval = CanEvaluate(lhs, rhs, "+", new string[] { "+", "*", "||" }, 0);
            if (eval) evaluatable.Add(lhs);
        }

        return new Result(evaluatable.Sum(), true);
    }
}
