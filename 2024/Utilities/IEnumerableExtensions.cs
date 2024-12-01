using System.Collections;

namespace AdventOfCode.Utilities;

public static class IEnumerableExtensions
{
    public static void PrintLines<T>(this IEnumerable<T> a)
    {
        // if (typeof(IEnumerable).IsAssignableFrom(typeof(T)))
        // {
        //     foreach (var i in a)
        //     {
        //         i.PrintLines();
        //     }
        // }
        foreach (var i in a) {
            Console.WriteLine(i);
        }
    }
}
