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

        return new Result(safeLevels.Where(i => i ==  true).Count(), true);
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
        var safeLevels = CheckSafety2(reportsList, 3, 3, true);

        return new Result(safeLevels.Where(i => i ==  true).Count(), true);
    }

    private (bool, int) CheckSingleReport(List<int> report, int maxIncrease, int maxDecrease, bool retry)
    {
        List<int> diffResults = new();

        for (var i = 0; i < report.Count() - 1; i++)
        {
            int? lastResult = null;
            bool resultSame = false;
            var diff = report[i] - report[i+1];
            int result;

            if (diff > 0 && diff <= maxDecrease) {
                result = -1;
                if (lastResult == null)
                {
                    lastResult = -1;
                }
                else
                {
                    resultSame = lastResult == -1;
                    lastResult = -1;
                }
            }
            else if (diff < 0 && Math.Abs(diff) <= maxIncrease)
            {
                result = 1;
                if (lastResult == null)
                {
                    lastResult = 1;
                }
                else
                {
                    resultSame = lastResult == 1;
                    lastResult = 1;
                }
            }
            else
            {
                result = 0;
                if (lastResult == null)
                {
                    lastResult = 0;
                }
                else
                {
                    resultSame = lastResult == 0;
                    lastResult = 0;
                }
            }

            diffResults.Add(result);

            if (!resultSame && retry)
            {
                (var try1Safety, var v1) = CheckSingleReport(report.Where((v, idx) => idx != i).ToList(), maxIncrease, maxDecrease, false);
                (var try2Safety, var v2) = CheckSingleReport(report.Where((v, idx) => idx != i + 1).ToList(), maxIncrease, maxDecrease, false);
                if (try1Safety == true)
                {
                    return (try1Safety, v1);
                }
                else if (try2Safety == true)
                {
                    return (try2Safety, v2);
                }
            }
        }
        var dist = diffResults.Distinct();
        var safetyCheck = dist.Count() == 1 && !dist.Contains(0);
        return (safetyCheck, dist.First());
    }

    private List<bool> CheckSafety2(List<List<int>> reports, int maxIncrease, int maxDecrease, bool retry)
    {
        List<bool> safeReports = new();

        foreach (var report in reports)
        {
            safeReports.Add(CheckSingleReport(report, maxIncrease, maxDecrease, retry).Item1);
        }

        return safeReports;
    }
}
