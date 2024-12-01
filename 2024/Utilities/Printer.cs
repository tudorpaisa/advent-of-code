using AdventOfCode.Exercises;

namespace AdventOfCode.Utilities;

public static class Printer
{
    public static void PrintRunDetails(int day, int part, bool runTest)
    {
        Console.WriteLine("\n##################################");
        Console.WriteLine($"Day: {day}");
        Console.WriteLine($"Part: {part}");
        Console.WriteLine($"Test Cases: {runTest}");
    }

    public static void PrintResult(Result result)
    {
        Console.WriteLine($"---\nCompleted: {result.Completed.ToString()}");
        Console.WriteLine($"Answer: {result.Answer.ToString()}");
    }
}
