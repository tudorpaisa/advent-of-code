namespace AdventOfCode.Utilities;

public static class ListExtensions
{
    private static void AddMultiple<T>(this List<T> list, T val, int repetitions)
    {
        for (var i = 0; i < repetitions; i++) list.Add(val);
    }

}
