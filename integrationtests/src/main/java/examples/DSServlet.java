package examples;

import java.io.IOException;
import java.io.PrintWriter;
import java.sql.Connection;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.sql.Statement;

import javax.annotation.Resource;
import javax.servlet.ServletException;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import javax.sql.DataSource;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

/**
 * Servlet implementation class DSServlet
 */
public class DSServlet extends HttpServlet {
    private static final long serialVersionUID = 1L;

    private static Logger logger = LogManager.getLogger(DSServlet.class);

    @Resource(name = "jdbc/DatasourceOne")
    DataSource ds1;

    public DSServlet() {
        super();
    }

    /**
     * @see HttpServlet#doGet(HttpServletRequest request, HttpServletResponse response)
     */
    @Override
    protected void doGet(HttpServletRequest request, HttpServletResponse response) throws ServletException, IOException {
        PrintWriter pw = response.getWriter();
        Connection connection1 = null;
        Statement selectStatement1 = null;
        boolean hasWarnings = false;
        boolean successful = true;

        final String NEWLINE = System.lineSeparator();
        StringBuilder sb = new StringBuilder();
        sb.append(NEWLINE);

        logger.info("Starting data source application test");
        try {
            try {
                sb.append("Working with data source \"jdbc/DatasourceOne\"").append(NEWLINE);
                connection1 = ds1.getConnection();
                selectStatement1 = connection1.createStatement();
                ResultSet rs1 = selectStatement1.executeQuery("SELECT * FROM pg_catalog.pg_tables");
                rs1.next();

                sb.append("Trying to get table names from table \"pg_tables\"").append(NEWLINE);
                String value1 = rs1.getString(1);
                sb.append("Result: ").append(value1).append(NEWLINE);
                if (value1 == null || value1.equals("")) {
                    successful = false;
                }
            } catch (SQLException e) {
                logger.error("Exception while working with data source \"jdbc/DatasourceOne\":", e);
                sb.append("Exception while working with data source \"jdbc/DatasourceOne\": " + e.getMessage()).append(NEWLINE);
                hasWarnings = true;
            } finally {
                if (selectStatement1 != null) {
                    try {
                        selectStatement1.close();
                    } catch (Exception e) {
                        logger.error("Exception while closing statement for data source \"jdbc/DatasourceOne\":", e);
                        sb.append("Exception while closing statement for data source \"jdbc/DatasourceOne\": " + e.getMessage()).append(NEWLINE);
                    }
                }
                if (connection1 != null) {
                    try {
                        connection1.close();
                    } catch (Exception e) {
                        logger.error("Exception while closing connection for data source \"jdbc/DatasourceOne\":", e);
                        sb.append("Exception while closing connection for data source \"jdbc/DatasourceOne\": " + e.getMessage()).append(NEWLINE);
                    }
                }
            }

            sb.append(NEWLINE);

            if (successful) {
                if (!hasWarnings) {
                    pw.write("OK");
                    System.out.println("OK");
                } else {
                    pw.write("OK WITH WARNINGS");
                    System.out.println("OK WITH WARNINGS");
                }
            } else {
                pw.write(sb.toString());
                System.out.println("FAILED");
            }

            logger.info("Test execution results:");
            logger.info(sb.toString());

            System.out.println(sb.toString());

            pw.flush();
        } finally {
            logger.info("Finishing data source application test");
            if (pw != null) {
                pw.close();
            }
        }
    }
}
