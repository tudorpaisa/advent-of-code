using AdventOfCode.Utilities;

namespace AdventOfCode.Exercises;

public class Exercise1 : IExercise
{
    public int GetDay() => 1;

    public void ExecutePart1(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        // input.PrintLines();

        List<List<int>> elvesFoods = [];

        var lineNo = 0;
        List<int> collection = [];
        List<int> sums = [];
        while (lineNo < input.Length)
        {
            var line = input[lineNo];
            if (line == String.Empty)
            {
                elvesFoods.Add(collection);
                sums.Add(collection.Aggregate((a, b) => a + b));
                collection = [];
            }
            else
            {
                Int32.TryParse(line, out var calories);
                collection.Add(calories);
            }
            lineNo++;
        }
        elvesFoods.Add(collection);
        sums.Add(collection.Aggregate((a, b) => a + b));
        collection = [];

        Helpers.PrintResult(sums.Max());
    }

    public void ExecutePart2(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        // input.PrintLines();

        List<List<int>> elvesFoods = [];

        var lineNo = 0;
        List<int> collection = [];
        List<int> sums = [];
        while (lineNo < input.Length)
        {
            var line = input[lineNo];
            if (line == String.Empty)
            {
                elvesFoods.Add(collection);
                sums.Add(collection.Aggregate((a, b) => a + b));
                collection = [];
            }
            else
            {
                Int32.TryParse(line, out var calories);
                collection.Add(calories);
            }
            lineNo++;
        }

        elvesFoods.Add(collection);
        sums.Add(collection.Aggregate((a, b) => a + b));
        collection = [];

        sums.Sort();
        sums.PrintLines();
        var top3 = sums.TakeLast(3);
        top3.PrintLines();

        Helpers.PrintResult(top3.Aggregate((a, b) => a + b));
    }
}
