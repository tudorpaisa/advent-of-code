using AdventOfCode.Utilities;

namespace AdventOfCode.Exercises;

public class Exercise2 : IExercise
{
    public int GetDay() => 2;

    public Result ExecutePart1(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);

        List<List<int>> reportsList = ParseLevels(input);
        var safeLevels = CheckSafety(reportsList, 3, 3);

        return new Result(safeLevels.Where(i => i ==  true).Count(), false);
    }

    private List<bool> CheckSafety(List<List<int>> reports, int maxIncrease, int maxDecrease)
    {
        List<bool> safeReports = new();

        foreach (var report in reports)
        {
            var maxIncDiff = 0;
            var maxDecDiff = 0;

            for (var i = 0; i < report.Count() - 1; i++)
            {
                var diff = report[i] - report[i+1];
                if (diff == 0)
                {
                    maxDecDiff = 99;
                    maxIncDiff = 99;
                }
                else if (diff > 0)
                {
                    maxDecDiff = Math.Max(maxDecDiff, Math.Abs(diff));
                }
                else
                {
                    maxIncDiff = Math.Max(maxIncDiff, Math.Abs(diff));
                }
            }
            // Console.WriteLine($"{maxIncDiff}, {maxDecDiff}");

            if (maxIncDiff > 0 && maxDecDiff > 0)
            {
                safeReports.Add(false);
            }
            else if (maxIncDiff == 0 && maxDecDiff > 0 && maxDecDiff <= maxDecrease)
            {
                safeReports.Add(true);
            }
            else if (maxDecDiff == 0 && maxIncDiff > 0 && maxIncDiff <= maxIncrease)
            {
                safeReports.Add(true);
            }
            else
            {
                safeReports.Add(false);
            }
        }

        return safeReports;
    }

    private List<List<int>> ParseLevels(string[] input)
    {
        List<List<int>> reportsList = new();

        foreach (var rawReport in input)
        {
            List<int> report = new();
            foreach (var rawLevel in rawReport.Split(" "))
            {
                int.TryParse(rawLevel, out var level);
                report.Add(level);
            }
            reportsList.Add(report);
        }
        return reportsList;
    }

    public Result ExecutePart2(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);

        List<List<int>> reportsList = ParseLevels(input);
        var safeLevels = CheckSafety2(reportsList, 3, 3);

        return new Result(safeLevels.Where(i => i ==  true).Count(), false);

        return new Result("", false);
    }

    private List<bool> CheckSafety2(List<List<int>> reports, int maxIncrease, int maxDecrease)
    {
        List<bool> safeReports = new();

        foreach (var report in reports)
        {
            var maxIncDiff = 0;
            var maxDecDiff = 0;

            for (var i = 0; i < report.Count() - 1; i++)
            {
                var diff = report[i] - report[i+1];
                if (diff == 0)
                {
                    maxDecDiff = 99;
                    maxIncDiff = 99;
                }
                else if (diff > 0)
                {
                    maxDecDiff = Math.Max(maxDecDiff, Math.Abs(diff));
                }
                else
                {
                    maxIncDiff = Math.Max(maxIncDiff, Math.Abs(diff));
                }
            }
            // Console.WriteLine($"{maxIncDiff}, {maxDecDiff}");

            if (maxIncDiff > 0 && maxDecDiff > 0)
            {
                safeReports.Add(false);
            }
            else if (maxIncDiff == 0 && maxDecDiff > 0 && maxDecDiff <= maxDecrease)
            {
                safeReports.Add(true);
            }
            else if (maxDecDiff == 0 && maxIncDiff > 0 && maxIncDiff <= maxIncrease)
            {
                safeReports.Add(true);
            }
            else
            {
                safeReports.Add(false);
            }
        }

        return safeReports;
    }
}
