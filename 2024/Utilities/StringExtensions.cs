namespace AdventOfCode.Utilities;

public static class StringExtensions
{
    public static void PrintLines(this string[] a)
    {
        foreach (var i in a) { Console.WriteLine(i); }
    }
}
